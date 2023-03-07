package chat

import "go.uber.org/zap"

type Settings struct {
	Logger *zap.Logger
}

type Service struct {
	logger *zap.Logger
}

func NewChatService(settings *Settings) *Service {
	return &Service{logger: settings.Logger}
}

func (srv *Service) CreateChat() {}
