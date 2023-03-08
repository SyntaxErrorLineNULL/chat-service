package repository

import "github.com/SyntaxErrorLineNULL/chat-service/repository/chat"

type Repository struct {
	Chat chat.IChatRepository
}

func New() *Repository {
	return &Repository{
		// create container loader
		Chat: chat.NewDefaultChatRepository(nil, nil),
	}
}
