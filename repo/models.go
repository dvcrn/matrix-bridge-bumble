package repo

import (
	"database/sql"
	"fmt"
	"time"

	"maunium.net/go/mautrix/bridge/bridgeconfig"

	matrix "github.com/dvcrn/matrix-bridgekit/matrix"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"maunium.net/go/mautrix/id"
)

type BumbleUserMatchType int

type SyncedMessage struct {
	BumbleID      string `gorm:"primarykey"`
	RoomID        id.RoomID
	MatrixEventID id.EventID `gorm:"index"`

	FromBumbleID string
	ToBumbleID   string

	Timestamp int64
	Content   string

	CreatedAt time.Time
}

type User struct {
	ID               int       `gorm:"primarykey"`
	MXID             id.UserID `gorm:"index"`
	SpaceMXID        id.RoomID
	ManagementRoomID id.RoomID

	RemoteID        string
	RemoteName      string
	DisplayName     string
	PermissionLevel bridgeconfig.PermissionLevel

	// Bumble Stuff
	// Bumble cookie stuff
	HDXUserID              string
	Aid                    string
	DeviceID               string
	Session                string
	FirstWebVisitID        string
	LastReferredWebVisitID string
	Domain                 string

	InitCompletedAt sql.NullTime
}

func (dbuser *User) ToMatrix() *matrix.User {
	return &matrix.User{
		MXID:             dbuser.MXID,
		ManagementRoomID: dbuser.ManagementRoomID,
		RemoteID:         dbuser.RemoteID,
		RemoteName:       dbuser.RemoteName,
		DisplayName:      dbuser.DisplayName,
		PermissionLevel:  dbuser.PermissionLevel,
	}
}

const (
	BumbleUserMatchTypeMyTurn BumbleUserMatchType = iota
	BumbleUserMatchTypeTheirTurn
	BumbleUserMatchTypeFullMatch
	BumbleUserMatchTypeDeleted
	BumbleUserMatchTypeExpired
	BumbleUserMatchTypeUnmatched
)

type BumbleUser struct {
	BumbleID      string    `gorm:"primarykey"`
	UserMXID      id.UserID `gorm:"index:unique"`
	RoomID        id.RoomID `gorm:"index:unique"`
	OwnerUserMXID id.UserID `gorm:"index"`

	Name                   string
	Age                    int
	Gender                 int
	ProfilePhotoID         string
	ProfilePhotoURL        string
	ProfilePhotoPreviewURL string
	MatchType              BumbleUserMatchType
	DisplayMessage         string

	ConnectionExpiredTimestamp int
	IsDeleted                  bool

	LastFullMessageSync      *time.Time
	LastProfileSync          *time.Time
	ProfilePhotoPreviewMXURL id.ContentURIString
	ProfilePhotoMXURL        id.ContentURIString
	RoomPhotoMXURL           id.ContentURIString

	RoomLockedAt *time.Time

	// AI Stuff
	AIWelcomeMessageEventID id.EventID
	AIWelcomeMessageSentAt  *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *BumbleUser) IsExpired() bool {
	return b.ConnectionExpiredTimestamp > 0
}

func (b *BumbleUser) IsUnmatched() bool {
	return b.MatchType == BumbleUserMatchTypeUnmatched
}

func (b *BumbleUser) IsUserDeleted() bool {
	return b.MatchType == BumbleUserMatchTypeDeleted || b.IsDeleted
}

func (b *BumbleUser) IsChatAvailable() bool {
	return !b.IsUserDeleted() && !b.IsExpired() && !b.IsUnmatched()
}

func (u *BumbleUser) NormalizeDisplayName() string {
	return fmt.Sprintf("%s 🐝", u.Name)
}

func (b *BumbleUser) ToMatrixGhost() *matrix.Ghost {
	return &matrix.Ghost{
		MXID:        b.UserMXID,
		RemoteID:    b.BumbleID,
		DisplayName: b.Name,
		UserName:    b.Name,
		AvatarURL:   id.ContentURI{},
	}
}

func (b *BumbleUser) ToMatrixRoom() *matrix.Room {
	return &matrix.Room{
		RemotedID:   b.BumbleID,
		MXID:        b.RoomID,
		Name:        b.Name,
		Encrypted:   false,
		PrivateChat: true,
		Ghosts:      []*matrix.Ghost{b.ToMatrixGhost()},
	}
}

func OpenConnection(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&BumbleUser{}); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&SyncedMessage{}); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	return db, nil
}
func GetBumbleUsersWithUnprocessedAIMessage(db *gorm.DB) ([]*BumbleUser, error) {
	var users []*BumbleUser
	tx := db.Where("match_type != ? AND ai_welcome_message_event_id IS NOT NULL AND ai_welcome_message_sent_at IS NULL", BumbleUserMatchTypeUnmatched).Order("updated_at DESC").Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return users, nil
}

func GetBumbleUserByBumbleID(db *gorm.DB, bumbleID string) (*BumbleUser, error) {
	var user BumbleUser
	tx := db.First(&user, "bumble_id = ?", bumbleID) // find product with integer primary key
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func GetBumbleUserByMXID(db *gorm.DB, userID id.UserID) (*BumbleUser, error) {
	var user BumbleUser
	tx := db.First(&user, "user_mx_id = ?", userID) // find product with integer primary key
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func GetRecentlyActiveBumbleUsers(db *gorm.DB, count int) ([]*BumbleUser, error) {
	var users []*BumbleUser
	tx := db.Where("match_type = ? OR match_type = ? OR match_type = ?", BumbleUserMatchTypeFullMatch, BumbleUserMatchTypeMyTurn, BumbleUserMatchTypeTheirTurn).Order("updated_at DESC").Limit(count).Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}

func GetActiveBumbleUsersByLatestProfileSyncOlderThan(db *gorm.DB, count int, olderThan time.Time) ([]*BumbleUser, error) {
	var users []*BumbleUser
	tx := db.Where("((match_type = ? OR match_type = ? OR match_type = ?) AND last_profile_sync < ?) OR last_profile_sync IS NULL", BumbleUserMatchTypeFullMatch, BumbleUserMatchTypeMyTurn, BumbleUserMatchTypeTheirTurn, olderThan).Order("last_profile_sync ASC").Limit(count).Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return users, nil
}

func GetAllBumbleUsers(db *gorm.DB) ([]BumbleUser, error) {
	var users []BumbleUser
	tx := db.Find(&users).Order("updated_at DESC")
	if tx.Error != nil {
		return nil, tx.Error
	}

	return users, nil
}

func GetBumbleUserByRoomID(db *gorm.DB, roomID id.RoomID) (*BumbleUser, error) {
	var user BumbleUser
	tx := db.First(&user, "room_id = ?", roomID) // find product with integer primary key
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func GetBumbleUserWithoutBackfill(db *gorm.DB) (*BumbleUser, error) {
	var user BumbleUser
	tx := db.First(&user, "last_full_message_sync is null") // find product with integer primary key
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func GetBumbleUserWithoutRoomPhoto(db *gorm.DB) (*BumbleUser, error) {
	var user BumbleUser
	tx := db.First(&user, "room_photo_mx_url == '' AND profile_photo_mx_url != ''") // find product with integer primary key
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func UpdateBumbleUser(db *gorm.DB, user *BumbleUser) error {
	user.UpdatedAt = time.Now()
	tx := db.Updates(user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func CreateBumbleUser(db *gorm.DB, user *BumbleUser) error {
	user.CreatedAt = time.Now()
	tx := db.Create(user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func GetSyncedMessage(db *gorm.DB, bumbleID string) (*SyncedMessage, error) {
	var user SyncedMessage
	tx := db.First(&user, "bumble_id = ?", bumbleID) // find product with integer primary key
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func GetSyncedMessageByBumbleIDOrEventID(db *gorm.DB, bumbleID string, evtID id.EventID) (*SyncedMessage, error) {
	var user SyncedMessage
	tx := db.First(&user, "matrix_event_id = ? OR bumble_id = ?", evtID, bumbleID) // find product with integer primary key
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func GetSyncedMessageByEventID(db *gorm.DB, evtID id.EventID) (*SyncedMessage, error) {
	var user SyncedMessage
	tx := db.First(&user, "matrix_event_id = ?", evtID) // find product with integer primary key
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func CreateSyncedMessageWithEventID(db *gorm.DB, bumbleID string, roomID id.RoomID, evtID id.EventID, fromID string, toID string, ts int64, msg string) (*SyncedMessage, error) {
	sm := &SyncedMessage{
		BumbleID:      bumbleID,
		RoomID:        roomID,
		FromBumbleID:  fromID,
		MatrixEventID: evtID,
		ToBumbleID:    toID,
		CreatedAt:     time.Now(),
		Timestamp:     ts,
		Content:       msg,
	}

	tx := db.Create(sm)

	return sm, tx.Error
}
func CreateSyncedMessage(db *gorm.DB, bumbleID string, roomID id.RoomID, fromID string, toID string, ts int64, msg string) (*SyncedMessage, error) {
	sm := &SyncedMessage{
		BumbleID:     bumbleID,
		RoomID:       roomID,
		FromBumbleID: fromID,
		ToBumbleID:   toID,
		CreatedAt:    time.Now(),
		Timestamp:    ts,
		Content:      msg,
	}

	tx := db.Create(sm)

	return sm, tx.Error
}

func UpdateSyncedMessage(db *gorm.DB, sm *SyncedMessage) error {
	tx := db.Updates(sm)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func GetAllUsers(db *gorm.DB) ([]*User, error) {
	var users []*User
	tx := db.Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}

func GetUserByMXID(db *gorm.DB, userID id.UserID) (*User, error) {
	var user User
	tx := db.First(&user, "mx_id = ?", userID)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func UpdateUser(db *gorm.DB, user *User) error {
	tx := db.Updates(user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func CreateUser(db *gorm.DB, user *User) error {
	tx := db.Create(user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func GetOldestMessageInRoom(db *gorm.DB, roomID id.RoomID) (*SyncedMessage, error) {
	var message SyncedMessage
	tx := db.Order("timestamp asc").First(&message, "room_id =?", roomID)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &message, nil
}

func GetNewestMessageInRoom(db *gorm.DB, roomID id.RoomID) (*SyncedMessage, error) {
	var message SyncedMessage
	tx := db.Order("timestamp desc").First(&message, "room_id =?", roomID)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &message, nil
}
