package bumble

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	bridgekit "github.com/dvcrn/matrix-bridgekit/bridgekit"
	matrix "github.com/dvcrn/matrix-bridgekit/matrix"
	"github.com/dvcrn/matrix-bumble/config"
	"github.com/dvcrn/matrix-bumble/repo"
	"gorm.io/gorm"
	"maunium.net/go/mautrix/bridge/commands"
)

type WrappedCommandEvent struct {
	*commands.Event

	Connector *BumbleConnector
	Bridge    *bridgekit.BridgeKit[*config.Config]
	User      *matrix.User
	Portal    *matrix.Room

	DB *gorm.DB
}

func extractCookie(curlString string) string {
	start := strings.Index(curlString, "Cookie:")
	if start == -1 {
		return ""
	}
	end := strings.Index(curlString[start:], "'")
	if end == -1 {
		return curlString[start:]
	}

	cookieString := curlString[start : start+end]
	return strings.Replace(cookieString, "Cookie: ", "", 1)
}

func extractDomain(curlString string) string {
	// Find the URL in the curl string
	urlStart := strings.Index(curlString, "https://")
	if urlStart == -1 {
		urlStart = strings.Index(curlString, "http://")
	}
	if urlStart == -1 {
		return ""
	}

	// Extract the URL
	urlEnd := strings.Index(curlString[urlStart:], "'")
	if urlEnd == -1 {
		urlEnd = strings.Index(curlString[urlStart:], "\"")
	}
	if urlEnd == -1 {
		return ""
	}
	fullURL := curlString[urlStart : urlStart+urlEnd]

	// Parse the URL
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return ""
	}

	// Extract and return the domain
	return parsedURL.Hostname()
}

func parseCookieString(cookieString string) map[string]string {
	cookies := map[string]string{}
	pairs := strings.Split(cookieString, "; ")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			cookies[parts[0]] = parts[1]
		}
	}
	return cookies
}

func fnLogin(ctx context.Context, ce *WrappedCommandEvent) {
	ce.Reply("okkkkk I try to log you in yah")
	ce.Reply("please paste the entire curl string taken from bumble website")

	ce.User.SetCommandState(&commands.CommandState{
		Action: "Login",
		Next: commands.MinimalHandlerFunc(func(ev *commands.Event) {
			ev.Reply("just a moment....")
			ev.Redact()

			if !strings.HasPrefix(strings.TrimSpace(ev.RawArgs), "curl") {
				ce.Reply("you need to paste the curl string, pls try again")
				return
			}

			cookieString := extractCookie(ev.RawArgs)
			cookies := parseCookieString(cookieString)
			domain := extractDomain(ev.RawArgs)

			// try to find user
			user, err := repo.GetUserByMXID(ce.DB, ce.User.MXID)
			isNewUser := errors.Is(err, gorm.ErrRecordNotFound)
			if err != nil {
				fmt.Println("User not found, creating...")
				user = &repo.User{
					MXID:            ce.User.MXID,
					PermissionLevel: ce.Bridge.Config.BridgeConfig.Permissions.Get(ce.User.MXID),
				}
			}
			user.Domain = domain

			// try to access cookies: aid, HDR-X-User-id, session, device_id, first_web_visit_id, last_referred_web_visit_id
			keys := []string{
				"aid",
				"HDR-X-User-id",
				"session",
				"device_id",
				"first_web_visit_id",
				//"last_referred_web_visit_id",
			}

			for _, key := range keys {
				if _, ok := cookies[key]; !ok {
					ev.Reply("could not find cookie: " + key)
					return
				}
			}

			for k, v := range cookies {
				switch k {
				case "aid":
					user.Aid = v
				case "HDR-X-User-id":
					user.HDXUserID = v
				case "session":
					user.Session = v
				case "device_id":
					user.DeviceID = v
				case "first_web_visit_id":
					user.FirstWebVisitID = v
				case "last_referred_web_visit_id":
					user.LastReferredWebVisitID = v
				}
			}

			if _, ok := cookies["aid"]; !ok {
				ev.Reply("aid not found in cookie string")
				return
			}

			bumbleClient := NewBumbleClient(domain, user.Aid, user.DeviceID, user.Session, user.HDXUserID, user.FirstWebVisitID, user.LastReferredWebVisitID)
			appStart, err := bumbleClient.FetchAppStart()
			if err != nil {
				ev.Reply("could not fetch app start: " + err.Error())
				return
			}

			if appStart.AppSettings == nil || appStart.AppSettings.BillingEmail == "" {
				ev.Reply("unable to authenticate, please try again")
				return
			}

			if len(appStart.ClientStartup.Host) > 0 {
				user.Domain = appStart.ClientStartup.Host[0]
			}

			// save
			if isNewUser {
				fmt.Println("Creating new user...")
				if err := repo.CreateUser(ce.DB, user); err != nil {
					ev.Reply("error creating user: " + err.Error())
					return
				}
			} else {
				fmt.Println("Updating user...")
				if err := repo.UpdateUser(ce.DB, user); err != nil {
					ev.Reply("error updating user: " + err.Error())
					return
				}
			}

			// set status
			ev.Reply("okkkkk I think I logged you in!!")
			ev.Reply("domain: " + domain)
			ev.Reply("bumble user ID: " + user.HDXUserID)
			ev.Reply("bumble email: " + appStart.AppSettings.BillingEmail)
			ev.Reply("device id: " + user.DeviceID)
			ce.Connector.StartUser(ctx, user)
		}),
	})
}

// func fnResyncAvatars(ctx context.Context, ce *WrappedCommandEvent) {
// 	userList, err := FetchUserListData(10, 0)
// 	if err != nil {
// 		ce.Reply(err.Error())
// 	}

// 	allUsers := userList.Section[0].Users
// 	if userList.Section[0].Name != "Conversations" {
// 		allUsers = userList.Section[1].Users
// 	}

// 	for _, user := range allUsers {
// 		dbUser, err := repo.GetBumbleUserByBumbleID(ce.DB, user.UserID)
// 		if err != nil {
// 			fmt.Println("err happened")
// 			continue
// 		}

// 		dbUser.ProfilePhotoID = user.ProfilePhoto.ID
// 		dbUser.ProfilePhotoURL = "https:" + user.ProfilePhoto.LargeURL
// 		dbUser.ProfilePhotoPreviewURL = "https:" + user.ProfilePhoto.PreviewURL

// 		if dbUser.ProfilePhotoID != "" {
// 			ghost := dbUser.ToMatrixGhost()
// 			ce.Bridge.GhostMaster.LoadGhost(ghost)
// 			contentURI, err := ce.Bridge.GhostMaster.UploadGhostAvatar(context.Background(), ghost, dbUser.ProfilePhotoPreviewURL)
// 			if err != nil {
// 				fmt.Println("couldn't upload avatar: ", err.Error())
// 				ce.Reply(fmt.Sprintf("Failed updating avatar of %s - %s", dbUser.Name, err.Error()))
// 				continue
// 			}

// 			dbUser.ProfilePhotoMXURL = contentURI.CUString()
// 			if err := repo.UpdateBumbleUser(ce.DB, dbUser); err != nil {
// 				ce.Reply("couldnt update user, but avatar changed: " + err.Error())
// 				return
// 			}

// 			ce.Reply(fmt.Sprintf("Updated %s avatar to %s", dbUser.Name, dbUser.ProfilePhotoPreviewURL))
// 		}

// 	}
// }

func (bum *BumbleConnector) GetCommands(ctx context.Context) []commands.Handler {
	var cmdLogin = &commands.FullHandler{
		Func: bum.wrapCommand(ctx, fnLogin),
		Name: "login",
		Help: commands.HelpMeta{
			Section:     commands.HelpSectionAuth,
			Description: "Login the bumble session",
			Args:        "[_bumble curl command_]",
		},
	}

	// var cmdResyncAvatars = &commands.FullHandler{
	// 	Func: bum.wrapCommand(ctx, fnResyncAvatars),
	// 	Name: "resync-avatars",
	// 	Help: commands.HelpMeta{
	// 		Section:     commands.HelpSectionAdmin,
	// 		Description: "Resyncs the last 10 avatars",
	// 	},
	// }

	var cmdMarkDumpRoomlist = &commands.FullHandler{
		Func: bum.wrapCommand(ctx, func(ctx context.Context, ce *WrappedCommandEvent) {
			users, err := repo.GetAllBumbleUsers(ce.DB)
			if err != nil {
				ce.Reply(err.Error())
				return
			}
			roomIds := []string{}
			for _, u := range users {
				roomIds = append(roomIds, u.RoomID.String())
			}

			msg := fmt.Sprintf("RoomIds: \n ```\n%s\n```", strings.Join(roomIds, ",\n"))
			ce.Reply(msg)
		}),
		Name: "dump-rooms",
		Help: commands.HelpMeta{
			Section:     commands.HelpSectionAdmin,
			Description: "Dumps all rooms",
		},
	}

	return []commands.Handler{cmdLogin, cmdMarkDumpRoomlist}
}
