package bumble

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dvcrn/matrix-bridgekit/bridgekit"
	"github.com/dvcrn/matrix-bridgekit/matrix"
	"github.com/dvcrn/matrix-bumble/ai"
	"github.com/dvcrn/matrix-bumble/bumble/bumbletypes"
	"github.com/dvcrn/matrix-bumble/config"
	"github.com/dvcrn/matrix-bumble/repo"
	"gorm.io/gorm"
	"maunium.net/go/mautrix/bridge"
	"maunium.net/go/mautrix/bridge/commands"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

var _ bridgekit.MatrixRoomEventHandler = (*BumbleConnector)(nil)
var _ bridgekit.BridgeConnector = (*BumbleConnector)(nil)
var _ bridgekit.MatrixRoomEventHandler = (*BumbleConnector)(nil)

type state struct {
	Bumble      *BumbleClient
	UserProfile *bumbletypes.UserExtended
	DBUser      *repo.User
}

type BumbleConnector struct {
	kit        *bridgekit.BridgeKit[*config.Config]
	db         *MiniDB
	gdb        *gorm.DB
	queue      *Queue
	cache      *Cache
	myBumbleID string

	ai *ai.AI

	matrixUsers  map[id.UserID]*matrix.User
	userState    map[id.UserID]*state
	roomSyncLock map[string]*sync.Mutex
	userSyncLock map[string]*sync.Mutex
}

func (bum *BumbleConnector) HandleMatrixMarkEncrypted(ctx context.Context, room *matrix.Room) error {
	//TODO implement me
	panic("implement me")
}

type MiniDB struct {
	Rooms  map[id.RoomID]*matrix.Room
	Ghosts map[id.UserID]*matrix.Ghost

	BackfilledMessages map[string]bool
}

func (bum *BumbleConnector) lockUser(room id.UserID) {
	if _, ok := bum.userSyncLock[room.String()]; !ok {
		bum.userSyncLock[room.String()] = &sync.Mutex{}
	}

	fmt.Printf("[user %s LOCK]\n", room.String())
	bum.userSyncLock[room.String()].Lock()
}

func (bum *BumbleConnector) unlockUser(room id.UserID) {
	fmt.Printf("[user %s UNLOCK]\n", room.String())
	if lock, ok := bum.userSyncLock[room.String()]; ok {
		lock.Unlock()
	}
}

func (bum *BumbleConnector) lockRoom(room *matrix.Room, reason string) {
	fmt.Println("attempting to lock room: ", room.MXID, "reason=", reason)
	if _, ok := bum.roomSyncLock[room.MXID.String()]; !ok {
		bum.roomSyncLock[room.MXID.String()] = &sync.Mutex{}
	}

	fmt.Printf("[room %s LOCK]\n", room.MXID)
	bum.roomSyncLock[room.MXID.String()].Lock()
	fmt.Println("room lock ok! ", room.MXID, "reason=", reason)
}

func (bum *BumbleConnector) unlockRoom(room *matrix.Room) {
	fmt.Printf("[room %s UNLOCK]\n", room.MXID)
	if lock, ok := bum.roomSyncLock[room.MXID.String()]; ok {
		lock.Unlock()
	}
}

func (bum *BumbleConnector) handleIncomingMatrixMessageEvent(room *matrix.Room, bridgeUser bridge.User, evt *event.Event) error {
	fmt.Println("got message event")
	// TODO: for debug. remove me.
	func(v interface{}) {
		j, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		buf := bytes.NewBuffer(j)
		fmt.Printf("%v\n", buf.String())
	}(evt)

	found, err := repo.GetSyncedMessageByEventID(bum.gdb, evt.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("error getting synced message: ", err)
		return err
	}

	if found != nil {
		fmt.Println("already found this msg in db, going to skip it")
		return nil
	}

	if strings.Contains(evt.Sender.String(), bum.kit.Bot.Localpart) {
		fmt.Println("message came from a ghost, just ignoring this")
		return nil
	}

	foundBumbleUser, err := repo.GetBumbleUserByRoomID(bum.gdb, room.MXID)
	if err != nil {
		fmt.Println("error getting user: ", err)
		return err
	}

	userState, ok := bum.userState[bridgeUser.GetMXID()]
	if !ok {
		return errors.New("could not get user state")
	}

	toBumbleUserID := foundBumbleUser.BumbleID
	fromBumbleUserID := bum.myBumbleID

	content, ok := evt.Content.Parsed.(*event.MessageEventContent)
	if !ok {
		return nil
	}

	msg := content.Body

	fmt.Println("msg:::")
	fmt.Println(toBumbleUserID)
	fmt.Println(fromBumbleUserID)
	fmt.Println(msg)

	sentmsg, err := userState.Bumble.sendMessage(fromBumbleUserID, toBumbleUserID, msg)
	fmt.Println(err)
	if err != nil {
		fmt.Println("error sending message: ", err)
		return err
	}

	// store it in synced messages, so it will not be considered as new
	if _, err := repo.CreateSyncedMessageWithEventID(bum.gdb, sentmsg.ChatMessage.UID, room.MXID, evt.ID, sentmsg.ChatMessage.FromPersonID, sentmsg.ChatMessage.ToPersonID, int64(sentmsg.ChatMessage.DateCreated*1000), sentmsg.ChatMessage.Mssg); err != nil {
		fmt.Println("err creating new synced msg ", err.Error())
	}

	return nil
}

func (bum *BumbleConnector) HandleMatrixRoomEvent(ctx context.Context, room *matrix.Room, user bridge.User, evt *event.Event) error {
	// lock the room so nothing is happening until we're done
	bum.lockRoom(room, "handling of matrix room event")
	defer bum.unlockRoom(room)

	switch evt.Type {
	case event.EventMessage:
		err := bum.handleIncomingMatrixMessageEvent(room, user, evt)
		if err != nil {
			fmt.Println("error handling incoming message: ", err)
			bum.kit.ReplyErrorMessage(ctx, evt, room, err)
			return err
		}
	default:
		fmt.Println("got unknown event type: ", evt.Type.String())
	}

	if err := bum.kit.MarkBotRead(ctx, room, evt); err != nil {
		fmt.Println("error replying read: ", err)
	}
	return nil
}

func (bum *BumbleConnector) wrapCommand(ctx context.Context, handler func(context.Context, *WrappedCommandEvent)) func(*commands.Event) {
	return func(ce *commands.Event) {
		user := ce.User.(*matrix.User)
		var portal *matrix.Room
		if ce.Portal != nil {
			portal = ce.Portal.(*matrix.Room)
		}

		handler(ctx, &WrappedCommandEvent{
			Event:     ce,
			Connector: bum,
			Bridge:    bum.kit,
			User:      user,
			Portal:    portal,
			DB:        bum.gdb,
		})
	}
}

func (bum *BumbleConnector) Init(ctx context.Context) error {
	//TODO implement me
	env, err := Process()
	if err != nil {
		return err
	}

	fmt.Println("[BumbleConnector Init] ", env)
	// TODO: for debug. remove me.
	func(v interface{}) {
		j, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		buf := bytes.NewBuffer(j)
		fmt.Printf("%v\n", buf.String())
	}(env)

	cache := NewCache()
	cache.Start(ctx)

	bum.queue = NewQueue(3)
	bum.cache = cache
	bum.db = &MiniDB{
		Rooms:              make(map[id.RoomID]*matrix.Room),
		Ghosts:             make(map[id.UserID]*matrix.Ghost),
		BackfilledMessages: make(map[string]bool),
	}
	bum.roomSyncLock = make(map[string]*sync.Mutex)

	for _, cmd := range bum.GetCommands(ctx) {
		bum.kit.RegisterCommand(cmd)
	}

	db, err := repo.OpenConnection(env.DatabasePath)
	if err != nil {
		fmt.Printf("error opening database connection: %v\n", err)
		return err
	}

	bum.gdb = db
	return nil
}

func (bum *BumbleConnector) UploadGhostProfilePicture(ghost *matrix.Ghost, room *matrix.Room, url string) (id.ContentURI, error) {
	fmt.Println("[SyncGhostProfilePicture] ", ghost.DisplayName, url)

	contentUrl, err := bum.kit.GhostMaster.UploadGhostAvatar(context.Background(), ghost, url)
	if err != nil {
		fmt.Println("failed to update avatar: ", err)
		return id.ContentURI{}, err
	}
	return contentUrl, nil
}

func (bum *BumbleConnector) SyncUserUpdate(ctx context.Context, u *matrix.User, foundUser *repo.BumbleUser, userListEntry *bumbletypes.User) {
	fmt.Println("syncing user update for user ", foundUser.Name)

	foundUserMatrix := foundUser.ToMatrixRoom()
	foundUserGhost := foundUser.ToMatrixGhost()

	now := time.Now()
	foundUser.LastProfileSync = &now
	defer func() {
		fmt.Println("defer -- updating user ")
		// TODO: for debug. remove me.
		func(v interface{}) {
			j, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			buf := bytes.NewBuffer(j)
			fmt.Printf("%v\n", buf.String())
		}(foundUser)

		if err := repo.UpdateBumbleUser(bum.gdb, foundUser); err != nil {
			fmt.Println("error updating user: ", err)
		} else {
			fmt.Println("successfully updated user")
		}
	}()

	if foundUser.DisplayMessage != userListEntry.DisplayMessage && !strings.Contains(userListEntry.DisplayMessage, "You connected") && !strings.Contains(userListEntry.DisplayMessage, "Conversation expired") {
		fmt.Println("display message changed, updating.", "existing", foundUser.DisplayMessage, "new", userListEntry.DisplayMessage)
		foundUser.DisplayMessage = userListEntry.DisplayMessage

		if err := repo.UpdateBumbleUser(bum.gdb, foundUser); err != nil {
			fmt.Println("error updating user: ", err)
			return
		} else {
			fmt.Println("successfully updated user display message to: ", foundUser.DisplayMessage)
		}

		// we'll download the latest 30 messages if the user status has changed
		bum.queue.Enqueue("latest_messages_"+userListEntry.UserID, func() { bum.SyncLatestRoomMessages(ctx, u, userListEntry.UserID, 30, true) })
	}

	if !foundUser.IsChatAvailable() {
		fmt.Println("user is not in an available state, skipping: ", foundUser.MatchType)
		return
	}

	if foundUser.ConnectionExpiredTimestamp != userListEntry.ConnectionExpiredTimestamp || foundUser.RoomLockedAt == nil && userListEntry.ConnectionExpiredTimestamp > 0 {
		fmt.Println("looks like this user is now disabled")
		foundUser.ConnectionExpiredTimestamp = userListEntry.ConnectionExpiredTimestamp
		foundUser.MatchType = repo.BumbleUserMatchTypeExpired
		t := time.Now()
		foundUser.RoomLockedAt = &t
		if err := repo.UpdateBumbleUser(bum.gdb, foundUser); err != nil {
			fmt.Println("error updating user: ", err)
			return
		}

		// mark chat as dead
		msg := format.RenderMarkdown("The connection with this user expired.", true, false)
		bum.kit.SendTimestampedBotMessageInRoom(ctx, foundUserMatrix, &msg, int64(foundUser.ConnectionExpiredTimestamp*1000))
		if err := bum.LockUserRoom(ctx, foundUser); err != nil {
			fmt.Println("error locking room: ", err)
		}
	}

	if (foundUser.IsDeleted != userListEntry.IsDeleted) || (foundUser.RoomLockedAt == nil && userListEntry.IsDeleted) {
		fmt.Println("looks like this user is now deleted")
		foundUser.IsDeleted = userListEntry.IsDeleted
		foundUser.MatchType = repo.BumbleUserMatchTypeDeleted
		t := time.Now()
		foundUser.RoomLockedAt = &t

		if err := repo.UpdateBumbleUser(bum.gdb, foundUser); err != nil {
			fmt.Println("error updating user: ", err)
			return
		}

		// mark chat as dead
		if err := bum.LockUserRoom(ctx, foundUser); err != nil {
			fmt.Println("error locking room: ", err)
		}

		msg := format.RenderMarkdown("The user has deleted their account", true, false)
		bum.kit.SendBotMessageInRoom(ctx, foundUserMatrix, &msg)
		return
	}

	if !userListEntry.IsMatch && foundUser.MatchType != repo.BumbleUserMatchTypeUnmatched {
		fmt.Println("looks like this user is no longer a match")
		fmt.Println("isMatch: ", userListEntry.IsMatch)
		fmt.Println("matchtype: ", foundUser.MatchType, "is unmatched? ", foundUser.MatchType == repo.BumbleUserMatchTypeUnmatched)

		foundUser.IsDeleted = userListEntry.IsDeleted
		foundUser.MatchType = repo.BumbleUserMatchTypeUnmatched
		t := time.Now()
		foundUser.RoomLockedAt = &t

		if err := repo.UpdateBumbleUser(bum.gdb, foundUser); err != nil {
			fmt.Println("error updating user: ", err)
			return
		}

		fmt.Println("NEW Found User MatchType: ", foundUser.MatchType, "is unmatched? ", foundUser.MatchType == repo.BumbleUserMatchTypeUnmatched)

		// mark chat as dead
		if err := bum.LockUserRoom(ctx, foundUser); err != nil {
			fmt.Println("error locking room: ", err)
		}

		msg := format.RenderMarkdown("This user is no longer a match 💔", true, false)
		bum.kit.SendBotMessageInRoom(ctx, foundUserMatrix, &msg)
		return
	}

	// photo URL changed, we should update that
	if foundUser.ProfilePhotoID != userListEntry.ProfilePhoto.ID && userListEntry.ProfilePhoto.PreviewURL != "" {
		fmt.Println("profile photo changed, updating: ", foundUser.ProfilePhotoID, " -> ", userListEntry.ProfilePhoto.ID)

		foundUser.ProfilePhotoURL = "https:" + userListEntry.ProfilePhoto.LargeURL
		foundUser.ProfilePhotoPreviewURL = "https:" + userListEntry.ProfilePhoto.PreviewURL
		foundUser.ProfilePhotoID = userListEntry.ProfilePhoto.ID

		if _, err := bum.UploadGhostProfilePicture(foundUserGhost, foundUserMatrix, foundUser.ProfilePhotoPreviewURL); err != nil {
			fmt.Println("failed to sync profile picture: ", err)
		}

		if err := repo.UpdateBumbleUser(bum.gdb, foundUser); err != nil {
			fmt.Println("error updating user: ", err)
			return
		}

		bum.SyncRoomAvatar(ctx, foundUser, time.Now().UnixMilli())
	}

	if userListEntry.ReplyTimeLeft == nil && userListEntry.PreMatchTimeLeft == nil {
		fmt.Println("Both ReplyTimeLeft and PreMatchTimeLeft are nil, this is probably now a match")
		if err := bum.UnlockUserRoom(ctx, foundUser); err != nil {
			fmt.Println("error unlocking room: ", err)
		}

		foundUser.MatchType = repo.BumbleUserMatchTypeFullMatch
	}

	if userListEntry.ReplyTimeLeft != nil {
		fmt.Println("reply time left changed, updating")
		if userListEntry.ReplyTimeLeft.Progress > 0 {
			if err := bum.UnlockUserRoom(ctx, foundUser); err != nil {
				fmt.Println("error unlocking room: ", err)
			}
		}
	}

}

func (bum *BumbleConnector) CreateMatrixUserFromBumble(ctx context.Context, u *matrix.User, bumbleUser *bumbletypes.User) (*repo.BumbleUser, error) {
	mt := repo.BumbleUserMatchTypeFullMatch
	if bumbleUser.IsDeleted {
		mt = repo.BumbleUserMatchTypeDeleted
	} else if bumbleUser.ConnectionExpiredTimestamp > 0 {
		mt = repo.BumbleUserMatchTypeExpired
	}

	foundUser, err := repo.GetUserByMXID(bum.gdb, u.MXID)
	if err != nil {
		return nil, err
	}

	user := repo.BumbleUser{
		BumbleID: bumbleUser.UserID,
		//UserMXID:                   "",
		//RoomID:                     "",
		OwnerUserMXID:          u.MXID,
		Name:                   bumbleUser.Name,
		Age:                    bumbleUser.Age,
		Gender:                 bumbleUser.Gender,
		ProfilePhotoURL:        "https:" + bumbleUser.ProfilePhoto.LargeURL,
		ProfilePhotoPreviewURL: "https:" + bumbleUser.ProfilePhoto.PreviewURL,
		ProfilePhotoID:         bumbleUser.ProfilePhoto.ID,
		DisplayMessage:         bumbleUser.DisplayMessage,

		IsDeleted:                  bumbleUser.IsDeleted,
		ConnectionExpiredTimestamp: bumbleUser.ConnectionExpiredTimestamp,
		MatchType:                  mt,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	fmt.Println("-- Creating ghost now....")
	ghost := bum.kit.GhostMaster.NewGhost(user.BumbleID, bumbleUser.NormalizeDisplayName(), bumbleUser.NormalizeUsername(), id.ContentURI{})
	if err := bum.kit.GhostMaster.AsGhost(ghost).Register(ctx); err != nil {
		fmt.Println("error ensuring as registered", err)
		return nil, err
	}
	bum.db.Ghosts[ghost.MXID] = ghost

	avatarContentUri, err := bum.kit.GhostMaster.UploadGhostAvatar(ctx, ghost, fmt.Sprintf("https:%s", bumbleUser.ProfilePhoto.PreviewURL))
	if err != nil {
		fmt.Println("failed to sync profile picture: ", err)
		avatarContentUri = id.ContentURI{}
	}

	user.UserMXID = ghost.MXID
	user.ProfilePhotoPreviewMXURL = avatarContentUri.CUString()

	room := matrix.NewRoom(bumbleUser.NormalizeDisplayName(), fmt.Sprintf(
		"%s on Bumble 🐝",
		bumbleUser.Name,
	), bum.kit.Bot, ghost)

	r, _, err := bum.kit.CreateRoom(ctx, room, u, avatarContentUri)
	if err != nil {
		fmt.Println("error creating room", err)
		return nil, err
	}

	bum.kit.GhostMaster.UpdateGhostName(ctx, ghost, bumbleUser.NormalizeDisplayName())
	bum.kit.RoomManager.AddRoomToUserSpace(ctx, foundUser.SpaceMXID, room)

	user.RoomID = r.MXID
	bum.db.Rooms[r.MXID] = r

	if err := repo.CreateBumbleUser(bum.gdb, &user); err != nil {
		fmt.Println("error creating user: ", err)
		return nil, err
	}

	bum.queue.Enqueue("sync_full_pic_"+ghost.DisplayName, func() {
		// bum.lockRoom(r)
		bum.lockRoom(r, "locking before sync full pic")
		defer bum.unlockRoom(r)

			fullImageContentUri, err := bum.uploadImageFromURL(ctx, fmt.Sprintf("https:%s", bumbleUser.ProfilePhoto.LargeURL))
			if err != nil {
				fmt.Println("failed to sync profile picture: ", err)
				return
			}

		user.ProfilePhotoMXURL = fullImageContentUri.CUString()
		if err := repo.UpdateBumbleUser(bum.gdb, &user); err != nil {
			fmt.Println("error updating user: ", err)
			return
		}
	})

	fmt.Println("Created user: ", user.Name)
	return &user, nil
}

func (bum *BumbleConnector) SyncUserList(ctx context.Context, user *matrix.User, maxLimit int) {
	fmt.Println("[SyncUserList] ", user.DisplayName)

	userState, ok := bum.userState[user.MXID]
	if !ok {
		return
	}

	fetchLimit := 30
	if maxLimit < 30 {
		fetchLimit = maxLimit
	}

	processedCount := 0
	for processedCount <= maxLimit {
		userList, err := userState.Bumble.FetchUserListData(fetchLimit, processedCount)
		if err != nil {
			fmt.Printf("error fetching user list: %v\n", err)
			return
		}

		if userList == nil || len(userList.Section) == 0 {
			fmt.Println("no user list data")
			return
		}

		allUsers := userList.Section[0].Users
		// if userList.Section[0].Name != "Conversations" {
		// 	allUsers = userList.Section[1].Users
		// }

		allUsers = append(allUsers, userList.Section[1].Users...)
		if len(allUsers) == 0 {
			return
		}

		for _, userListEntry := range allUsers {
			fmt.Println("got user list entry: ", userListEntry.Name, " -- ", userListEntry.IsMatch)
			processedCount += 1

			// check if we already have
			// if we already do, we gotta update the user state
			foundUser, err := repo.GetBumbleUserByBumbleID(bum.gdb, userListEntry.UserID)
			if err == nil && foundUser != nil {
				// TODO REMOVE MME
				// if processedCount > 6 && processedCount < 10 {
				// 	bum.queue.Enqueue("gen_ai_welcome_msg"+user.DisplayName, func() {
				// 		bum.GenerateAIWelcomeMessage(ctx, user, foundUser)
				// 	})
				// }

				bum.SyncUserUpdate(ctx, user, foundUser, &userListEntry)
				continue
			}

			if !errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Print("already have user, or different error. skipping\n", err.Error())
				continue
			}

			fmt.Println("dont have user yet, creating ghost + room")
			createdUser, err := bum.CreateMatrixUserFromBumble(ctx, user, &userListEntry)
			if err != nil {
				fmt.Println("error creating user: ", err)
				continue
			}

			if userListEntry.PreMatchTimeLeft != nil && userListEntry.ReplyTimeLeft != nil {
				fmt.Println("is a prematch")

				if userListEntry.ReplyTimeLeft.Progress == 0 {
					fmt.Println("progress is '0', assuming this means we can't reply yet so marking the room as readonly")
					if err := bum.LockUserRoom(ctx, createdUser); err != nil {
						fmt.Println("error locking room: ", err)
					}
				}
			}

			// send initial match message

			msgs, err := bum.SyncLatestRoomMessages(ctx, user, createdUser.BumbleID, 100, false)
			if err != nil {
				fmt.Println("backfilling failed: ", err.Error())
			}

			bum.lockRoom(createdUser.ToMatrixRoom(), "locking to send welcome message")
			if len(msgs) > 0 {
				latestMsg := msgs[len(msgs)-1]

				content := format.RenderMarkdown("You've matched with "+userListEntry.Name+" 🍯 (but you'll need to wait until they texted you)", true, false)
					if msg, err := bum.kit.SendTimestampedBotMessageInRoom(ctx, createdUser.ToMatrixRoom(), &content, latestMsg.Timestamp); err != nil {
						fmt.Println("error sending message: ", err)
					} else {
						_ = bum.kit.GhostMaster.AsUserGhost(ctx, user).MarkRead(ctx, createdUser.ToMatrixRoom().MXID, msg.EventID)
					}
				} else {
					content := format.RenderMarkdown("You've matched with "+userListEntry.Name+" 🍯 (but you'll need to wait until they texted you)", true, false)
					if msg, err := bum.kit.SendBotMessageInRoom(ctx, createdUser.ToMatrixRoom(), &content); err != nil {
						fmt.Println("error sending message: ", err)
					} else {
						_ = bum.kit.GhostMaster.AsUserGhost(ctx, user).MarkRead(ctx, createdUser.ToMatrixRoom().MXID, msg.EventID)
					}
				}

			fmt.Println("---- trying to send photo")
			fmt.Println("----- ", userListEntry.ProfilePhoto)

			content2 := &event.MessageEventContent{
				MsgType: event.MsgImage,
				URL:     createdUser.ProfilePhotoMXURL,
				Info: &event.FileInfo{
					ThumbnailURL: createdUser.ProfilePhotoMXURL,
					MimeType:     "image/jpeg",
				},
				FileName: fmt.Sprintf("%s_avatar.jpg", createdUser.BumbleID),
				BeeperGalleryImages: []*event.MessageEventContent{
					{
						Info: &event.FileInfo{
							ThumbnailURL: createdUser.ProfilePhotoMXURL,
							MimeType:     "image/jpeg",
						},
						URL: createdUser.ProfilePhotoMXURL,
					},
				},
			}
				if msg, err := bum.kit.SendBotMessageInRoom(ctx, createdUser.ToMatrixRoom(), content2); err != nil {
					fmt.Println("error sending message: ", err)
				} else {
					_ = bum.kit.GhostMaster.AsUserGhost(ctx, user).MarkRead(ctx, createdUser.ToMatrixRoom().MXID, msg.EventID)
				}

			bum.unlockRoom(createdUser.ToMatrixRoom())
			//if userState.DBUser.InitCompletedAt.Valid && !userState.DBUser.InitCompletedAt.Time.IsZero() {
			//	return
			//	aiPlaceholderMessage, err := bum.kit.SendBotTextMessageInRoom(ctx, createdUser.ToMatrixRoom(), "Loading AI... 🔄")
			//	if err != nil {
			//		fmt.Println("error sending message: ", err)
			//		return
			//	}
			//
			//	createdUser.AIWelcomeMessageEventID = aiPlaceholderMessage.EventID
			//	repo.UpdateBumbleUser(bum.gdb, createdUser)
			//
			//	bum.queue.Enqueue("gen_ai_welcome_msg"+createdUser.Name, func() {
			//		bum.GenerateAIWelcomeMessage(ctx, user, createdUser)
			//	})
			//}

			if processedCount >= maxLimit {
				break
			}
		}
	}
}

func (bum *BumbleConnector) handleIncomingBumbleMessage(ctx context.Context, u *matrix.User, msg *bumbletypes.ChatMessage) {
	fmt.Println("got chat message")

	userState, ok := bum.userState[u.MXID]
	if !ok {
		fmt.Println("could not find userstate for user ", u.DisplayName)
		return
	}

	// check rooms
	user, err := repo.GetBumbleUserByBumbleID(bum.gdb, msg.ToPersonID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user, err = repo.GetBumbleUserByBumbleID(bum.gdb, msg.FromPersonID)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("user for message not found, attempting to create -- ", err)
		openChat, err := bum.cachedClientOpenChat(ctx, userState.Bumble, msg.FromPersonID)
		if err != nil {
			fmt.Println("tried to open chat but failed: ", err)
			return
		}

		created, err := bum.CreateMatrixUserFromBumble(ctx, u, openChat.ChatUser)
		if err != nil {
			fmt.Println("error creating user: ", err)
			return
		}

		fmt.Println("done creating new user. ", created.Name, "has been created.")
		user = created
	}

	room := user.ToMatrixRoom()
	bum.kit.RoomManager.LoadRoom(room)

	fmt.Println("found room: ", room.Name)

	// lock the room so nothing is happening until we're done
	bum.lockRoom(room, "lock to sync latest room messages")
	defer bum.unlockRoom(room)

	if _, err := repo.GetSyncedMessage(bum.gdb, msg.UID); err == nil {
		fmt.Println("already synced")
		return
	}

	if user.LastFullMessageSync == nil {
		bum.unlockRoom(room)
		fmt.Println("this is a new room, going to try to prefill first...")
		if _, err := bum.SyncLatestRoomMessages(ctx, u, user.BumbleID, 300, true); err != nil {
			fmt.Println("error syncing latest messages: ", err)
		}
		bum.lockRoom(room, "lock for getting synced messages")
		// no need to unlock here because defer call will still hit

		if _, err := repo.GetSyncedMessage(bum.gdb, msg.UID); err == nil {
			fmt.Println("already synced")
			return
		}
	}

	syncdMsg, err := repo.CreateSyncedMessage(bum.gdb, msg.UID, room.MXID, msg.FromPersonID, msg.ToPersonID, int64(msg.DateCreated*1000), msg.Mssg)
	if err != nil {
		fmt.Println("err creating new synced msg ", err.Error())
	}

	// if room ID = sender ID, message came from the other person
	// otherwise it's us
	content := format.HTMLToContent(msg.Mssg)
	sender := room.MainIntent()

	ts := int64(msg.DateCreated * 1000)
	bum.SyncRoomAvatar(ctx, user, ts)
	bum.SyncRoomDisplayName(ctx, user, ts)

	if room.RemotedID != msg.FromPersonID {
		fmt.Println("going to send message from USER in room now")
		res, err := bum.kit.SendUserMessageInRoom(ctx, room, u, &content)
		if err != nil {
			fmt.Println("error sending message: ", err)
			return
		}

		syncdMsg.MatrixEventID = res.EventID
		if err := repo.UpdateSyncedMessage(bum.gdb, syncdMsg); err != nil {
			fmt.Println("error updating synced message: ", err)
			return
		}

		return
	}

	fmt.Println("going to send message in room now")
	res, err := bum.kit.SendMessageInRoom(ctx, room, sender, &content)
	if err != nil {
		fmt.Println("error sending message: ", err)
		return
	}

	syncdMsg.MatrixEventID = res.EventID
	if err := repo.UpdateSyncedMessage(bum.gdb, syncdMsg); err != nil {
		fmt.Println("error updating synced message: ", err)
		return
	}
}

func (bum *BumbleConnector) StartWebsocket(ctx context.Context, user *matrix.User, path string, sequence string, userUid string, deviceID string) {
	userState, ok := bum.userState[user.MXID]
	if !ok {
		fmt.Println("could not find userstate")
		return
	}

	go userState.Bumble.OpenWebsocket(ctx, path, sequence, userUid, deviceID, func(msg *bumbletypes.EventMessage) {
		fmt.Println("got msg: ")
		// TODO: for debug. remove me.
		func(v interface{}) {
			j, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			buf := bytes.NewBuffer(j)
			fmt.Printf("%v\n", buf.String())
		}(msg)

		// on message:
		// 1. check if we have room already
		// 2. create room if we don't

		if msg.Body[0].ChatMessage != nil {
			bum.handleIncomingBumbleMessage(ctx, user, msg.Body[0].ChatMessage)
		}

		if msg.Body[0].PersonNotice != nil {
			bum.SyncUserList(ctx, user, 10)
			bum.SlowlyCheckRecentConnections(ctx, user, 5)
		}
	}, true)
}

func (bum *BumbleConnector) StartUser(ctx context.Context, dbUser *repo.User) {
	if dbUser.Aid == "" || dbUser.DeviceID == "" || dbUser.HDXUserID == "" {
		fmt.Println("no aid or device id, skipping", dbUser.DisplayName)
		return
	}

	domain := "bumble.com"
	if dbUser.Domain != "" {
		domain = dbUser.Domain
	}

	bumbleClient := NewBumbleClient(domain, dbUser.Aid, dbUser.DeviceID, dbUser.Session, dbUser.HDXUserID, dbUser.FirstWebVisitID, dbUser.LastReferredWebVisitID)

	user := dbUser.ToMatrix()
	bum.matrixUsers[user.MXID] = user

	fmt.Println("Starting user: ", user.DisplayName, user.MXID.String())

	c := bum.kit.AS.Client(user.MXID)
	newIntent, newAccessToken, err := bum.kit.Bridge.DoublePuppet.Setup(ctx, user.MXID, c.AccessToken, true)
	if err != nil {
		fmt.Println("Error setting up double puppet: ", err)
	} else {
		user.AccessToken = newAccessToken
		user.DoublePuppetIntent = newIntent
		fmt.Println("double puppetting ok")
	}

	appStart, err := bumbleClient.FetchAppStart()
	if err != nil {
		fmt.Println("error fetching app start: ", err)
		return
	}

	ownUser, err := bumbleClient.GetUser(dbUser.HDXUserID)
	if err != nil {
		fmt.Println("error fetching own user: ", err)
		return
	}

	bum.userState[user.MXID] = &state{
		Bumble:      bumbleClient,
		UserProfile: ownUser,
		DBUser:      dbUser,
	}

	if len(appStart.ClientStartup.Host) > 0 {
		bumbleClient.Domain = appStart.ClientStartup.Host[0]
		dbUser.Domain = appStart.ClientStartup.Host[0]
	}

	if appStart == nil || appStart.CometConfiguration == nil {
		fmt.Println("error: appStart or appStart cometconfiguration is nil")
		return
	}

	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	if err := bum.kit.GhostMaster.UpdateName(ctx, testGhost, "Spooky Spooky Ghost"); err != nil {
	// 		fmt.Println("err updating ghost name: ", err.Error())
	// 	}

	// 	content := format.RenderMarkdown("See? I can also update my own name", true, false)
	// 	bum.kit.SendMessageInRoom(ctx, createdRoom, createdRoom.MainIntent(), &content)
	// }()

	// content := format.RenderMarkdown("Hello, I'm a bot", true, false)
	// bum.kit.SendBotMessageInRoom(ctx, createdRoom, &content)

	// content = format.RenderMarkdown("Hello, I'm a ghost", true, false)
	// bum.kit.SendMessageInRoom(ctx, createdRoom, createdRoom.MainIntent(), &content)

	// continue

	if dbUser.SpaceMXID == "" {
		fmt.Println("user does not have a space yet, gonna create one")
		res, err := bum.kit.RoomManager.CreatePersonalSpace(
			context.Background(),
			user,
			"Bumble Chat",
			"Your personal space for Bumble conversations",
		)
		if err != nil {
			fmt.Println("error creating personal space: ", err)
			return
		}

		dbUser.SpaceMXID = res.RoomID
		if err := repo.UpdateUser(bum.gdb, dbUser); err != nil {
			fmt.Println("error updating user: ", err)
			return
		}
	} else {
		bum.kit.RoomManager.AddUserToRoom(ctx, dbUser.SpaceMXID, user)
	}

	if !dbUser.InitCompletedAt.Valid {
		bum.SyncUserList(ctx, user, 50)
		dbUser.InitCompletedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		if err := repo.UpdateUser(bum.gdb, dbUser); err != nil {
			fmt.Println("error updating user: ", err)
		}
	} else {
		bum.SyncUserList(ctx, user, 20)
	}

		// path := "<redacted>"
		// sequence := "<redacted>"
		// userUid := "<redacted>"
		// deviceID := "<redacted>"
	userUid := dbUser.Aid
	deviceID := dbUser.DeviceID
	cometPath := appStart.CometConfiguration.Path
	cometSeq := appStart.CometConfiguration.Sequence

	bum.StartWebsocket(ctx, user, cometPath, cometSeq, userUid, deviceID)

	fmt.Println("bumble connector started")

	//go bum.ReconcileAIWelcomeMessages(ctx, user)
	go bum.SlowlyBackfillRemainingConversations(ctx, user)
	// go func() {
	// 	bum.SlowlyCheckRecentConnections(ctx, user, 5)
	// }()
	// run sync userlist every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			<-ticker.C
			fmt.Println("TICK resyncing userlist")
			bum.SyncUserList(ctx, user, 10)
			// bum.SlowlyCheckRecentConnections(ctx, user, 5)
			//bum.ReconcileAIWelcomeMessages(ctx, user)
		}
	}()

	ticker15 := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			<-ticker15.C
			bum.SlowlyCheckNotUpdatedConnections(ctx, user, 10)
		}
	}()
}

func (bum *BumbleConnector) Start(ctx context.Context) {
	users, err := repo.GetAllUsers(bum.gdb)
	if err != nil {
		fmt.Println("error getting all users: ", err)
		return
	}

	for _, dbUser := range users {
		bum.StartUser(ctx, dbUser)
	}

}

func (bum *BumbleConnector) Stop() {
	fmt.Println("[Stop]")
	//TODO implement me
	// panic("implement me")
}

func (bum *BumbleConnector) SlowlyCheckRecentConnections(ctx context.Context, u *matrix.User, amount int) {
	userState, ok := bum.userState[u.MXID]
	if !ok {
		fmt.Println("could not get userstate")
		return
	}

	fmt.Println("Starting slowlyyy check recent conversations")
	time.Sleep(10 * time.Second)

	users, err := repo.GetRecentlyActiveBumbleUsers(bum.gdb, amount)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("no matches found")
		return
	}

	fmt.Println("SlowlyCheckRecentConnections: checking for users. amount: ", len(users))
	for _, user := range users {
		// check if context is cancelled
		select {
		case <-ctx.Done():
			fmt.Println("SlowlyCheckRecentConnections: context cancelled")
			return
		default:
			// continue
		}

		fmt.Println("checking user: ", user.BumbleID, user.Name)
		openChat, err := userState.Bumble.ClientOpenChat(user.BumbleID, 0)
		if err != nil {
			fmt.Printf("error fetching chat instance: %v\n", err)
			return
		}

		bum.SyncUserUpdate(ctx, u, user, openChat.ChatUser)
		time.Sleep(5 * time.Second)
	}
}

func (bum *BumbleConnector) SlowlyCheckNotUpdatedConnections(ctx context.Context, u *matrix.User, amount int) {
	userState, ok := bum.userState[u.MXID]
	if !ok {
		fmt.Println("could not get userstate")
		return
	}

	fmt.Println("Starting slowlyyy check recent conversations that are older than 5 days")
	time.Sleep(10 * time.Second)

	olderThan5days := time.Now().AddDate(0, 0, -5)
	users, err := repo.GetActiveBumbleUsersByLatestProfileSyncOlderThan(bum.gdb, amount, olderThan5days)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("no matches found")
		return
	}

	fmt.Println("SlowlyCheckNotUpdatedConnections: checking for users. amount: ", len(users))
	for _, user := range users {
		// check if context is cancelled
		select {
		case <-ctx.Done():
			fmt.Println("SlowlyCheckNotUpdatedConnections: context cancelled")
			return
		default:
			// continue
		}

		fmt.Println("checking user: ", user.BumbleID, user.Name)
		openChat, err := bum.cachedClientOpenChat(ctx, userState.Bumble, user.BumbleID)
		if err != nil {
			fmt.Printf("error fetching chat instance: %v\n", err)
			return
		}

		bum.SyncUserUpdate(ctx, u, user, openChat.ChatUser)
		time.Sleep(15 * time.Second)
	}
}

func (bum *BumbleConnector) SlowlyBackfillRemainingConversations(ctx context.Context, u *matrix.User) {
	fmt.Println("Starting slow backfill")
	time.Sleep(10 * time.Second)
	for {
		// check if context is cancelled
		select {
		case <-ctx.Done():
			fmt.Println("SlowlyBackfillRemainingConversations: context cancelled")
			return
		default:
			// continue
		}

		user, err := repo.GetBumbleUserWithoutBackfill(bum.gdb)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}

		fmt.Println("backfilling: ", user.BumbleID)
		_, err = bum.SyncLatestRoomMessages(ctx, u, user.BumbleID, 100, false)
		if err != nil {
			fmt.Println("backfilling failed: ", err.Error())
		}
	}
}

func (bum *BumbleConnector) SyncRoomDisplayName(ctx context.Context, bumUser *repo.BumbleUser, ts int64) error {
	fmt.Println("syncing room display name for: ", bumUser.Name)
	// check the display name
	foundUserGhost := bumUser.ToMatrixGhost()
	foundUserRoom := bumUser.ToMatrixRoom()

	displayName, err := bum.kit.Bot.GetDisplayName(ctx, bumUser.UserMXID)
	if err != nil {
		fmt.Println("error getting room display name: ", err)
	} else {
		if displayName.DisplayName != bumUser.NormalizeDisplayName() {
			fmt.Println("display name changed, updating")

			// update ghost name, update room name
			if err := bum.kit.GhostMaster.UpdateGhostName(ctx, foundUserGhost, bumUser.NormalizeDisplayName()); err != nil {
				fmt.Println("error updating ghost name: ", err)
			}

			if err := bum.kit.RoomManager.SetRoomName(ctx, foundUserRoom, bum.kit.GhostMaster.AsRoomGhost(foundUserRoom), bumUser.NormalizeDisplayName()); err != nil {
				fmt.Println("error updating room name: ", err)
			}
		}
	}

	return nil
}

func (bum *BumbleConnector) SyncRoomAvatar(ctx context.Context, bumUser *repo.BumbleUser, ts int64) error {
	if (bumUser.RoomPhotoMXURL == "" || bumUser.RoomPhotoMXURL != bumUser.ProfilePhotoPreviewMXURL) && bumUser.ProfilePhotoPreviewMXURL != "" {
		fmt.Println("Photo has changed, updating room")
		cu, err := bumUser.ProfilePhotoPreviewMXURL.Parse()
		if err != nil {
			fmt.Println("error parsing url: ", err)
			return err
		}

		mroom := bumUser.ToMatrixRoom()
		bum.kit.RoomManager.LoadRoom(mroom)

		fmt.Println("setting room avatar now")
		if err := bum.kit.RoomManager.InsertSetRoomAvatarEvent(ctx, mroom, mroom.BotIntent, cu, ts); err != nil {
			fmt.Println("error setting room avatar: ", err)
			return err
		}

		bumUser.RoomPhotoMXURL = bumUser.ProfilePhotoPreviewMXURL
		if err := repo.UpdateBumbleUser(bum.gdb, bumUser); err != nil {
			fmt.Println("error updating user: ", err)
		}
	}

	return nil
}

func (bum *BumbleConnector) SyncLatestRoomMessages(ctx context.Context, user *matrix.User, bumbleID string, limit int, notify bool) ([]*matrix.Message, error) {
	bumbleUser, err := repo.GetBumbleUserByBumbleID(bum.gdb, bumbleID)
	if err != nil {
		fmt.Println("error getting room: ", err)
		return nil, err
	}

	mroom := bumbleUser.ToMatrixRoom()
	bum.kit.RoomManager.LoadRoom(mroom)

	userState, ok := bum.userState[user.MXID]
	if !ok {
		fmt.Println("could not find userstate")
		return nil, errors.New("could not find userstate")
	}

	fmt.Println("found room in db: ", bumbleUser.Name)
	bum.lockRoom(mroom, "lock for SyncLatestRoomMessages")
	defer bum.unlockRoom(mroom)

	fmt.Println("starting to backfill...")

	res, err := userState.Bumble.ClientOpenChat(bumbleID, limit)
	if err != nil {
		fmt.Printf("error fetching chat instances: %v\n", err)
		return nil, err
	}

	if res == nil {
		fmt.Println("clientOpenChat returned nil")
		return nil, errors.New("clientOpenChat returned nil")
	}

	fmt.Println("got messages: ", len(res.ChatMessages))

	msgs := []*matrix.Message{}
	bmsgs := []*bumbletypes.ChatMessage{}
	for _, msg := range res.ChatMessages {
		_, err := repo.GetSyncedMessage(bum.gdb, msg.UID)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}

		sender := user.MXID
		receiver := bum.kit.GhostMaster.AsRoomGhost(mroom).UserID
		//receiver := mroom.MainIntent().UserID
		if mroom.RemotedID == msg.FromPersonID {
			//sender = mroom.MainIntent().UserID
			sender = bum.kit.GhostMaster.AsRoomGhost(mroom).UserID
			receiver = user.MXID
		}

		msgs = append(msgs, &matrix.Message{
			FromMXID:  sender,
			ToMXID:    receiver,
			RoomID:    mroom.MXID,
			Content:   format.HTMLToContent(msg.Mssg),
			Timestamp: int64(msg.DateCreated * 1000),
		})

		bmsgs = append(bmsgs, msg)
	}

	now := time.Now()
	if len(msgs) == 0 {
		fmt.Println("no messages to backfill")

		bumbleUser.LastFullMessageSync = &now
		if err := repo.UpdateBumbleUser(bum.gdb, bumbleUser); err != nil {
			fmt.Println("error updating user: ", err)
		}

		return []*matrix.Message{}, nil
	}

	latestMsg := msgs[len(msgs)-1]

	bum.SyncRoomAvatar(ctx, bumbleUser, latestMsg.Timestamp)

	if err := bum.kit.BackfillMessages(ctx, mroom, user, msgs, notify); err != nil {
		fmt.Printf("error backfilling messages: %v\n", err)
		return nil, err
	}

	// mark all messages as processed
	for _, bmsg := range bmsgs {
		if _, err := repo.CreateSyncedMessage(bum.gdb, bmsg.UID, mroom.MXID, bmsg.FromPersonID, bmsg.ToPersonID, int64(bmsg.DateCreated*1000), bmsg.Mssg); err != nil {
			fmt.Printf("error creating synced message: %v\n", err)
		}
	}

	bumbleUser.LastFullMessageSync = &now
	if err := repo.UpdateBumbleUser(bum.gdb, bumbleUser); err != nil {
		fmt.Println("error updating user: ", err)
	}

	return msgs, nil
}

func (bum *BumbleConnector) GetRoom(ctx context.Context, roomID id.RoomID) *matrix.Room {
	fmt.Println("[BumbleConnector GetRoom] ", roomID.String())

	r, err := repo.GetBumbleUserByRoomID(bum.gdb, roomID)
	if err != nil {
		fmt.Println("error getting room: ", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// bum.kit.AS.Client(bum.user.MXID).LeaveRoom(context.Background(), roomID)
			return &matrix.Room{
				MXID: roomID,
			}
		}

		return &matrix.Room{
			MXID: roomID,
		}
	}

	userState, ok := bum.userState[r.OwnerUserMXID]
	if !ok {
		fmt.Println("could not get owner")
		return &matrix.Room{
			MXID: roomID,
		}
	}

	fmt.Println("found room in db: ", r.Name)

	if !r.IsChatAvailable() {
		fmt.Println("room is not available")
		return &matrix.Room{
			MXID: roomID,
		}
	}

	mroom := r.ToMatrixRoom()
	bum.kit.RoomManager.LoadRoom(mroom)

	ownerUser, ok := bum.matrixUsers[r.OwnerUserMXID]
	if !ok {
		fmt.Println("could not get owner")
		return &matrix.Room{
			MXID: roomID,
		}
	}

	// if no full message sync done yet, enqueue syncing it
	if r.LastFullMessageSync == nil {
		bum.queue.Enqueue("latest_room_messages_"+roomID.String(), func() {
			_, err := bum.SyncLatestRoomMessages(ctx, ownerUser, r.BumbleID, 100, false)
			if err != nil {
				fmt.Println("error syncing room messages: ", err)
				return
			}

			openChat, err := bum.cachedClientOpenChat(ctx, userState.Bumble, r.BumbleID)
			if err != nil {
				fmt.Printf("error fetching chat instances: %v\n", err)
				return
			}
			bum.SyncUserUpdate(ctx, ownerUser, r, openChat.ChatUser)
		})

		return mroom
	} else {
		openChat, err := bum.cachedClientOpenChat(ctx, userState.Bumble, r.BumbleID)
		if err != nil {
			fmt.Printf("error fetching chat instances: %v\n", err)
			return nil
		}
		bum.SyncUserUpdate(ctx, ownerUser, r, openChat.ChatUser)
	}

	return mroom
}

func (bum *BumbleConnector) GetAllRooms(ctx context.Context) []bridge.Portal {
	fmt.Println("[GetAllIPortals]")
	//TODO implement me
	panic("implement me")
}

func (bum *BumbleConnector) GetUser(ctx context.Context, id id.UserID, create bool) *matrix.User {
	fmt.Println("[GetIUser] ", id.String(), " create ", create)

	if user, ok := bum.matrixUsers[id]; ok {
		return user
	}

	dbuser, err := repo.GetUserByMXID(bum.gdb, id)
	if err == nil {
		bum.matrixUsers[id] = dbuser.ToMatrix()
		return dbuser.ToMatrix()
	}

	fmt.Println("did not find user in db, will create one I guess. create flag: ", create)
	if !create {
		return nil
	}

	if err := repo.CreateUser(bum.gdb, &repo.User{
		MXID:            id,
		PermissionLevel: bum.kit.Config.BridgeConfig.Permissions.Get(id),
	}); err != nil {
		fmt.Println("error creating user: ", err)
		return nil
	}

	return &matrix.User{
		MXID: id,
	}
}

func (bum *BumbleConnector) IsGhost(ctx context.Context, userID id.UserID) bool {
	fmt.Println("[IsGhost] ", userID.String())
	return false
	//TODO implement me
	// panic("implement me")
}

func (bum *BumbleConnector) GetGhost(ctx context.Context, userID id.UserID) *matrix.Ghost {
	fmt.Println("[GetIGhost] ", userID.String())
	if g, ok := bum.db.Ghosts[userID]; ok {
		fmt.Println("found ghost in db: ", g.DisplayName)
		return g
	}

	return nil
}

func (bum *BumbleConnector) SetManagementRoom(ctx context.Context, user *matrix.User, roomID id.RoomID) {
	fmt.Println("[SetManagementRoom] ", user.MXID.String(), " ", roomID.String())

	foundUser, err := repo.GetUserByMXID(bum.gdb, user.MXID)
	if err != nil {
		fmt.Println("could not find user to set management room")
		return
	}

	foundUser.ManagementRoomID = roomID
	if err := repo.UpdateUser(bum.gdb, foundUser); err != nil {
		fmt.Println("error updating user: ", err)
	}

	bum.matrixUsers[user.MXID].ManagementRoomID = roomID
}

func NewBumbleConnector(kit *bridgekit.BridgeKit[*config.Config]) *BumbleConnector {
	claude := ai.NewClaude(os.Getenv("ANTHROPIC_API_KEY"))
	ai := ai.NewAi(claude)

	br := &BumbleConnector{
		ai:          ai,
		kit:         kit,
		matrixUsers: make(map[id.UserID]*matrix.User),
		userState:   make(map[id.UserID]*state),
	}

	return br
}
