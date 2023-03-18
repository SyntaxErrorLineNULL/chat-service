package user

import (
	"github.com/SyntaxErrorLineNULL/chat-service/repository"
	"go.uber.org/zap"
)

type Settings struct {
	DB     *repository.Repository
	Logger *zap.Logger
}

type Service struct {
	db     *repository.Repository
	logger *zap.Logger
}

func New(settings *Settings) *Service {
	return &Service{}
}
