package chat

import (
	"context"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
	"github.com/SyntaxErrorLineNULL/chat-service/repository"
	gonanoid "github.com/matoous/go-nanoid/v2"
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

func NewChatService(settings *Settings) *Service {
	return &Service{db: settings.DB, logger: settings.Logger}
}

func (srv *Service) CreateChat(ctx context.Context, input *domain.Chat) error {
	l := srv.logger.Sugar().With("Create")
	l.Debug("creating new chat")

	if input == nil {
		l.Error(zap.Error(ErrEmpty), "empty request data")
		return ErrEmpty
	}

	id, _ := gonanoid.New(20)
	input.ID = id

	err := srv.db.Chat.Create(ctx, input)
	if err != nil {
		l.Error(zap.Error(err), "failed create chat")
		return err
	}

	l.Info("successfully create chat")
	return nil
}
