package repository

import (
	"github.com/SyntaxErrorLineNULL/chat-service/repository/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Repository struct {
	User user.UserRepository
}

func New(client *mongo.Client, logger *zap.Logger) *Repository {
	return &Repository{
		// create container loader
		User: user.NewDefaultUserRepository(client, logger),
	}
}
