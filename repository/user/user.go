package user

import (
	"context"
	"errors"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

// ErrEmpty returns when request was with empty data
var ErrEmpty = errors.New("empty")

// ErrNotFound returns when record doesn't exist in database
var ErrNotFound = errors.New("not found")

// ErrCannotFind returns when request was with incorrect data to search user
var ErrCannotFind = errors.New("cannot find")

// ErrInternal returns when something went wrong in repository
var ErrInternal = errors.New("internal error")

// ErrInvalidArgument returns when request was with invalid data
var ErrInvalidArgument = errors.New("invalid argument")

type DefaultUserRepository struct {
	client *mongo.Client
	col    *mongo.Collection
	logger *zap.Logger
}

func NewDefaultUserRepository(client *mongo.Client, logger *zap.Logger) *DefaultUserRepository {
	return &DefaultUserRepository{
		client: client,
		col:    client.Database("chat-service").Collection("user"),
		logger: logger,
	}
}

func (r *DefaultUserRepository) Create(ctx context.Context, usr *domain.User) error {
	l := r.logger.Sugar().With("Create")
	start := time.Now()
	if usr == nil {
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "user is nil")
		return ErrInvalidArgument
	}

	// generate user id
	usr.ID = uuid.New().String()
	_, err := r.col.InsertOne(ctx, usr)
	if err != nil {
		l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "insert error")
		return ErrInternal
	}

	l.Info(zap.Duration("duration", time.Since(start)), "successful user created")
	return nil
}
