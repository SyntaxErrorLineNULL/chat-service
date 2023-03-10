package chat

import (
	"context"
	"errors"
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

func (srv *Service) Create(ctx context.Context, input *domain.Chat) error {
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

func (srv *Service) CreatePersonalChatWithUser(ctx context.Context, ch *domain.Chat, userID string) (*domain.Chat, error) {
	l := srv.logger.Sugar().With("CreatePersonalChatWithUser")
	l.Debug("creating personal chat with user")

	// TODO: create get user id use request

	_, _, err := srv.isPossibleCreateChat(ctx, userID)
	if err != nil {
		l.Error(zap.Error(err), "no ability to create chat")
		return nil, err
	}

	var chat *domain.Chat
	chat, err = srv.db.Chat.FindPersonalChatBetweenUsers(ctx, "", "")
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	} else {
		// return chat if it already exists
		return chat, nil
	}

	if ch.Title == "" {
		// TODO: if chat title empty from user names you need to return the name of the chat, something like: Chat between OwnerID and UserID
	}

	return nil, nil
}

func (srv *Service) isPossibleCreateChat(ctx context.Context, userID string) (*domain.User, *domain.User, error) {
	l := srv.logger.Sugar().With("isPossibleCreateChat")
	l.Debug("checking the ability to create a chat")

	// TODO: create gRPC request to user service
	return nil, nil, nil
}
