package service

import "github.com/SyntaxErrorLineNULL/chat-service/service/chat"

type Provider struct {
	Chat chat.IChat
}

func New() *Provider {
	return &Provider{Chat: chat.NewChatService(nil)}
}
