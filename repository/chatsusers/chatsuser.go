package chatsusers

import (
	"context"
	"errors"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
	"github.com/SyntaxErrorLineNULL/chat-service/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type DefaultChatsUsersRepository struct {
	client *mongo.Client
	col    *mongo.Collection
	logger *zap.Logger
}

func NewDefaultChatsUsersRepository(client *mongo.Client, logger *zap.Logger) *DefaultChatsUsersRepository {
	return &DefaultChatsUsersRepository{
		client: client,
		col:    client.Database("chat-service").Collection("chats_users"),
		logger: logger,
	}
}

// Create a structure to store the chat user
func (r *DefaultChatsUsersRepository) Create(ctx context.Context, chu *domain.ChatsUsers) error {
	l := r.logger.Sugar().With("Create")
	start := time.Now()
	if chu.ID == "" {
		l.Error(zap.Error(repository.ErrEmpty), zap.Duration("duration", time.Since(start)), "empty request")
		return repository.ErrEmpty
	}

	_, err := r.col.InsertOne(ctx, chu)
	if err != nil {
		l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "insert error")
		return err
	}

	return nil
}

// Find chats users by chat ID and user ID value
func (r *DefaultChatsUsersRepository) Find(ctx context.Context, chatID, uid string) (*domain.ChatsUsers, error) {
	l := r.logger.Sugar().With("Find")
	start := time.Now()
	if chatID == "" {
		l.Error(zap.Error(repository.ErrEmpty), zap.Duration("duration", time.Since(start)), "empty request")
		return nil, repository.ErrEmpty
	}

	ch := &domain.ChatsUsers{}
	filter := bson.D{{Key: "chat_id", Value: chatID}, {Key: "user_id", Value: uid}}
	err := r.col.FindOne(ctx, filter).Decode(ch)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "not found in database")
			return nil, repository.ErrNotFound
		}
		l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "find error")
		return nil, err
	}

	return ch, nil
}
