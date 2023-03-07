package chat

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type DefaultChatRepository struct {
	client *mongo.Client
	col    *mongo.Collection
	logger *zap.Logger
}

func NewDefaultChatRepository(client *mongo.Client, logger *zap.Logger) *DefaultChatRepository {
	return &DefaultChatRepository{
		client: client,
		col:    client.Database("chat-service").Collection("chat"),
		logger: logger,
	}
}
