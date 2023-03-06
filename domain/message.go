package domain

// Message information
type Message struct {
	ID         string `json:"id" bson:"id"`
	ChatID     string `json:"chat_id" bson:"chat_id"`
	FromID     string `json:"from_id" bson:"from_id"`
	CreateDate int64  `json:"create_date" bson:"create_date"`
	// message type, indicates the content in the message. Can take values: text, image, video, document.
	Type string `json:"type" bson:"type"`
	// a link to the multimedia content of the message. You can get the link by first uploading the content to our server.
	Media    string `json:"media" bson:"media"`
	Body     string `json:"body" bson:"body"`
	UpdateAt int64  `json:"update_at" bson:"update_at"`
	// previous versions of the message (message structure object)
	Viewed bool `json:"viewed" bson:"viewed"`
	//
	Reaction string `json:"reaction,omitempty" bson:"reaction"`
}

const (
	// MessageTypeText is default type of the message
	MessageTypeText = "text"
	// MessageTypeImage if the media type jpg,png or other image format type.
	MessageTypeImage = "image"
	// MessageTypeVideo if the video
	MessageTypeVideo = "video"
	// MessageTypeDocument if the media type document
	MessageTypeDocument = "document"
	// MessageTypeSticker if the sticker image or document
	MessageTypeSticker = "sticker"
)
