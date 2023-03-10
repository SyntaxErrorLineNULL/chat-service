package chatsusers

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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
