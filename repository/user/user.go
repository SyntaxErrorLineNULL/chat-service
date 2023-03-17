package user

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// ErrEmpty returns when request was with empty data
var ErrEmpty = errors.New("empty")

// ErrNotFound returns when record doesn't exist in database
var ErrNotFound = errors.New("not found")

// ErrCannotFind returns when request was with incorrect data to search user
var ErrCannotFind = errors.New("cannot find")

// ErrInternal returns when something went wrong in repository
var ErrInternal = errors.New("internal error")

type DefaultChatRepository struct {
	client *mongo.Client
	col    *mongo.Collection
	logger *zap.Logger
}

func NewDefaultUserRepository(client *mongo.Client, logger *zap.Logger) *DefaultChatRepository {
	return &DefaultChatRepository{
		client: client,
		col:    client.Database("chat-service").Collection("user"),
		logger: logger,
	}
}
