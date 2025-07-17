package model

import (
	"time"

	"github.com/google/uuid"
)

// User table gorm model
type UserGorm struct {
	// Base Fields
	ID           uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username     string     `json:"username" gorm:"not null;Index"`
	Email        string     `json:"email" gorm:"size:254;not null;unique"`
	PasswordHash string     `json:"password_hash" gorm:"not null"`
	FirstName    *string    `json:"first_name,omitempty"`
	LastName     *string    `json:"last_name,omitempty"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	LastOnline   *time.Time `json:"last_online,omitempty"`
	Bio          *string    `json:"bio,omitempty"`
	DateOfBirth  time.Time  `json:"date_of_birth" gorm:"not null"`
	Location     *string    `json:"location,omitempty"`
	IsVerified   bool       `json:"is_verified" gorm:"default:false"`

	// Foreign Key for Profile Picture
	ProfilePictureAssetID *uuid.UUID `json:"profile_picture_asset_id" gorm:"type:uuid"`

	// Relations (Has One / Has Many)
	ProfilePictureAsset    *MediaAssetGorm         `gorm:"foreignKey:ProfilePictureAssetID"` // Relation: A user has one profile picture
	SentMessages           []MessageGorm           `gorm:"foreignKey:AuthorID"`              // Relation: A user sends many messages
	UserFriendConnectsSent []UserFriendConnectGorm `gorm:"foreignKey:User1ID"`               // Relation: A user initiates many friend connections
	UserFriendConnectsRecv []UserFriendConnectGorm `gorm:"foreignKey:User2ID"`               // Relation: A user receives many friend connections
	OwnedServers           []ServerGorm            `gorm:"foreignKey:OwnerID"`               // Relation: A user owns many servers
	ServerUserConnects     []ServerUserConnectGorm `gorm:"foreignKey:UserID"`                // Relation: A user is connected to many servers
	UploadedMediaAssets    []MediaAssetGorm        `gorm:"foreignKey:UploadedByUserID"`      // Relation: A user uploads many media assets
}

func (*UserGorm) TableName() string {
	return "users"
}

// MediaAsset table gorm model
type MediaAssetGorm struct {
	// Base Fields
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UrlPath string    `gorm:"not null"` // URL path to the file

	// Foreign Key for Uploader (Optional)
	UploadedByUserID *uuid.UUID `gorm:"type:uuid"` // Optional: Track who uploaded it

	// Relations
	UploadedByUser      *UserGorm  `gorm:"foreignKey:UploadedByUserID"`      // Relation: A media asset can be uploaded by a user
	UserProfilePictures []UserGorm `gorm:"foreignKey:ProfilePictureAssetID"` // Relation: A media asset can be a profile picture for multiple users
}

func (*MediaAssetGorm) TableName() string {
	return "media_assets"
}

// UserFriendConnect table gorm model
type UserFriendConnectGorm struct {
	// Composite Primary Keys (Foreign Keys)
	User1ID uuid.UUID `gorm:"not null;type:uuid;primaryKey;autoIncrement:false"`
	User2ID uuid.UUID `gorm:"not null;type:uuid;primaryKey;autoIncrement:false"`

	// Base Fields
	Status      Status
	RequestedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	AcceptedAt  *time.Time

	// Relations
	User1 UserGorm `gorm:"foreignKey:User1ID"` // Relation: Connects to the first user in the friendship
	User2 UserGorm `gorm:"foreignKey:User2ID"` // Relation: Connects to the second user in the friendship
}

func (*UserFriendConnectGorm) TableName() string {
	return "user_friend_connects"
}

// MessageGorm table gorm model
type MessageGorm struct {
	// Base Fields
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AuthorID    uuid.UUID   `gorm:"not null;type:uuid"`
	MessageType MessageType `gorm:"not null"`
	Content     string      `gorm:"not null;index:idx_content_type_search,priority:1"`
	CreatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   *time.Time  `gorm:"autoUpdateTime"`

	// Foreign Key for Reply
	ReplyTo *uuid.UUID `gorm:"type:uuid"` // Can be null if not a reply

	// Relations
	Author         UserGorm      `gorm:"foreignKey:AuthorID"` // Relation: A message has one author
	ReplyToMessage *MessageGorm  `gorm:"foreignKey:ReplyTo"`  // Relation: A message can reply to another message
	Replies        []MessageGorm `gorm:"foreignKey:ReplyTo"`  // Relation: A message can have many replies
}

func (*MessageGorm) TableName() string {
	return "messages"
}

// ServerGorm table gorm model
type ServerGorm struct {
	// Base Fields
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title      string     `gorm:"index:idx_title_visibility_search,priority:1"`
	OwnerID    uuid.UUID  `gorm:"type:uuid"`
	CreatedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	Visibility Visibility `gorm:"index:idx_title_visibility_search,priority:2"`

	// Foreign Key for Icon
	IconAssetID *uuid.UUID `gorm:"type:uuid"`

	// Relations
	Owner UserGorm `gorm:"foreignKey:OwnerID"` // Relation: A server has one owner
}

func (*ServerGorm) TableName() string {
	return "servers"
}

// ServerUserConnectGorm table gorm model
type ServerUserConnectGorm struct {
	// Composite Primary Keys (Foreign Keys)
	ServerID uuid.UUID `gorm:"not null;type:uuid;primaryKey;autoIncrement:false"`
	UserID   uuid.UUID `gorm:"not null;type:uuid;primaryKey;autoIncrement:false"`

	// Relations
	Server ServerGorm `gorm:"foreignKey:ServerID"` // Relation: Connects to the server
	User   UserGorm   `gorm:"foreignKey:UserID"`   // Relation: Connects to the user
}

func (*ServerUserConnectGorm) TableName() string {
	return "server_user_connects"
}

// Enums for gorm models
type Status string

const (
	StatusPending  Status = "pending"
	StatusAccepted Status = "accepted"
	StatusDeclined Status = "declined"
	StatusBlocked  Status = "blocked"
)

type Visibility string

const (
	VisibilityPrivate Visibility = "private"
	VisibilityPublic  Visibility = "public"
)

type ChannelType string

const (
	ChannelTypeChat  ChannelType = "chat"
	ChannelTypeVoice ChannelType = "voice"
)

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeVideo MessageType = "video"
	MessageTypeAudio MessageType = "audio"
)

type AssetType string

const (
	AssetTypeImage    AssetType = "image"
	AssetTypeVideo    AssetType = "video"
	AssetTypeAudio    AssetType = "audio"
	AssetTypeDocument AssetType = "document"
)
