package domain

// Chat information
type Chat struct {
	ID                string   `json:"id" bson:"id"`
	Title             string   `json:"title" bson:"title"`
	CreateDate        int64    `json:"create-date" bson:"create_date"`
	Type              string   `json:"type" bson:"type"`
	Participants      []string `json:"participants" bson:"participants"`
	Deleted           bool     `json:"deleted" bson:"deleted"`
	OwnerID           string   `json:"owner_id" bson:"owner_id"`
	Unread            int64    `json:"unread" bson:"unread"`
	PinnedMessagesIDs []string `json:"pinned_messages,omitempty" bson:"pinned_messages"`
	Label             string   `json:"label,omitempty" bson:"label"`
	PinnedMessageID   string   `json:"pinned_message_id,omitempty" bson:"pinned_message_id"`
	LastMessage       *Message `json:"last_message,omitempty" bson:"last_message"`
}

// Description of the chat type.
var (
	// ChatTypePersonal is default type of the chat
	ChatTypePersonal = "personal"
	// ChatTypeGroup is a chat type of group
	ChatTypeGroup = "group"
	// ChatTypeArchive this chat type is intended for archiving chats
	ChatTypeArchive = "archive"
)
