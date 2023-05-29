package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
	"github.com/SyntaxErrorLineNULL/chat-service/repository"
	cht "github.com/SyntaxErrorLineNULL/chat-service/repository/chat"
	"github.com/SyntaxErrorLineNULL/chat-service/service/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ErrEmpty returns when request was with empty data
var ErrEmpty = errors.New("empty request")

type Settings struct {
	DB     *repository.Repository
	Logger *zap.Logger
	Usr    *user.Service
}

type Service struct {
	db     *repository.Repository
	logger *zap.Logger
	usr    *user.Service
}

func New(settings *Settings) *Service {
	l := settings.Logger.With(zap.Namespace("chat service instance"))
	return &Service{
		db:     settings.DB,
		logger: l,
	}
}

func (srv *Service) methodLogger(ctx context.Context, name string) *zap.Logger {
	return srv.logger.With(
		zap.String("layer: {user-service} ", fmt.Sprintf("method: [%s]", name)),
	)
}

func (srv *Service) Create(ctx context.Context, input *domain.Chat) (*domain.Chat, error) {
	// Get the method logger with additional context information.
	l := srv.methodLogger(ctx, "Create")
	l.Debug("creating new chat")

	if input == nil {
		l.Error("empty request data", zap.Error(ErrEmpty))
		return nil, ErrEmpty
	}

	input.ID = uuid.New().String()

	err := srv.db.Chat.Create(ctx, input)
	if err != nil {
		l.Error("failed create chat", zap.Error(err))
		return nil, err
	}

	res, err := srv.db.Chat.Find(ctx, input)
	if err != nil {
		l.Error("failed find chat", zap.Error(err))
		return nil, err
	}

	l.Info("successfully create chat")
	return res, nil
}

func (srv *Service) CreatePersonalChat(ctx context.Context, ch *domain.Chat, uid string, user *domain.User) (*domain.Chat, error) {
	// Get the method logger with additional context information.
	l := srv.methodLogger(ctx, "CreatePersonalChat")
	l.Debug("creating personal chat with other user")

	if uid == "" || user == nil {
		l.Error("uid or user is empty", zap.Error(ErrEmpty))
		return nil, ErrEmpty
	}

	own, err := srv.usr.Find(ctx, &domain.User{ID: uid})
	if err != nil {
		l.Error("failed find owner", zap.Error(err))
		return nil, err
	}

	member, err := srv.isPossibleCreateChat(ctx, user)
	if err != nil {
		l.Error("no ability to create chat", zap.Error(err))
		return nil, err
	}

	var chat *domain.Chat
	chat, err = srv.db.Chat.FindPersonalChatBetweenUsers(ctx, own.ID, member.ID)
	if err != nil {
		if !errors.Is(err, cht.ErrNotFound) {
			l.Error("failed find chat between participants")
			return nil, err
		}
	} else {
		// return chat if it already exists
		l.Info("chat between this users already exists")
		return chat, nil
	}

	if ch.Title == "" {
		// TODO: if chat title empty from user names you need to return the name of the chat, something like: Chat between OwnerID and UserID
		ch.Title = fmt.Sprintf("Chat between user %s and user %s", own.UserName, member.UserName)
	}

	// return new chat
	return srv.Create(ctx, &domain.Chat{
		Title:        ch.Title,
		Type:         domain.ChatTypePersonal,
		Participants: []string{uid, member.ID},
		Deleted:      false,
		OwnerID:      uid,
		LastMessage:  nil,
		Unread:       0,
	})
}

func (srv *Service) isPossibleCreateChat(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Get the method logger with additional context information.
	l := srv.methodLogger(ctx, "isPossibleCreateChat")
	l.Debug("checking the ability to create a chat")

	member, err := srv.usr.Find(ctx, user)
	if err != nil {
		l.Error("failed find user", zap.Error(err))
		return nil, err
	}

	return member, nil
}
