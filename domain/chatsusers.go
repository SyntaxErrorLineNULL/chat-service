package domain

// ChatsUsers this is a structure that describes information about the chat user.
type ChatsUsers struct {
	ID     string `json:"id" bson:"id"`
	ChatID string `json:"chat_id" bson:"chat_id"`
	// user ID, which corresponds to the user ID in the system.
	UserID string `json:"user_id" bson:"user_id"`
	// the time when the participant joins the chat.
	AddedAt int64 `json:"added_at" bson:"added_at"`
	// the message identifier from which the chat message history is available.
	StartMessageID int64 `json:"start_message_id,omitempty" bson:"start_message_id"`
	// message identifier up to which the chat message history is available.
	EndMessageID int64 `json:"end_message_id,omitempty" bson:"end_message_id"`
	// MaxReadDate time of the last message read
	MaxReadDate int64 `json:"max_read_date" bson:"max_read_date"`
}
