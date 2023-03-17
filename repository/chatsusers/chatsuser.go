package chatsusers

import (
	"context"
	"errors"
	"time"

	"github.com/SyntaxErrorLineNULL/chat-service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// ErrEmpty returns when request was with empty data
var ErrEmpty = errors.New("empty")

// ErrNotFound returns when record doesn't exist in database
var ErrNotFound = errors.New("not found")

// ErrInternal returns when something went wrong in repository
var ErrInternal = errors.New("internal error")

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
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "empty request")
		return ErrEmpty
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
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "empty request")
		return nil, ErrEmpty
	}

	ch := &domain.ChatsUsers{}
	filter := bson.D{{Key: "chat_id", Value: chatID}, {Key: "user_id", Value: uid}}
	err := r.col.FindOne(ctx, filter).Decode(ch)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "chats users record not found in database")
			return nil, ErrNotFound
		}
		l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "find error")
		return nil, err
	}

	return ch, nil
}

// Update change record
func (r *DefaultChatsUsersRepository) Update(ctx context.Context, chu *domain.ChatsUsers) error {
	l := r.logger.Sugar().With("Update")
	start := time.Now()
	if chu.ID == "" {
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "empty request")
		return ErrEmpty
	}

	update := bson.M{
		"$set": bson.M{
			"start_message_id": chu.StartMessageID,
			"end_message_id":   chu.EndMessageID,
			"max_read_date":    chu.MaxReadDate,
		},
	}

	_, err := r.col.UpdateOne(ctx, bson.M{"chat_id": chu.ChatID}, update, options.Update().SetUpsert(true))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "chats users record not found in database")
			return ErrNotFound
		}
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "update error")
		return ErrInternal
	}

	return nil
}

// Delete deleting a chat participant's record.
func (r *DefaultChatsUsersRepository) Delete(ctx context.Context, chatID, uid string) error {
	l := r.logger.Sugar().With("Update")
	start := time.Now()
	if chatID == "" || uid == "" {
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "empty request")
		return ErrEmpty
	}

	filter := bson.D{{Key: "chat_id", Value: chatID}, {Key: "user_id", Value: uid}}
	_, err := r.col.DeleteOne(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "chat user record not found in database")
			return ErrNotFound
		}
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "error delete chat user")
		return ErrInternal
	}

	return nil
}
