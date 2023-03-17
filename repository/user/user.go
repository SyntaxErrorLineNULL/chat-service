package user

import (
	"context"
	"errors"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
	"github.com/SyntaxErrorLineNULL/chat-service/repository/dto"
	"github.com/SyntaxErrorLineNULL/chat-service/repository/mapper"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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
	u := &mapper.UserMapper{}
	_, err := r.col.InsertOne(ctx, u.ToDTO(usr))
	if err != nil {
		l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "insert error")
		return ErrInternal
	}

	l.Info(zap.Duration("duration", time.Since(start)), "successful user created")
	return nil
}

func (r *DefaultUserRepository) Find(ctx context.Context, usr *domain.User) (*domain.User, error) {
	l := r.logger.Sugar().With("Find")
	start := time.Now()
	if usr == nil {
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "user is nil")
		return nil, ErrInvalidArgument
	}

	filled := false
	match := bson.A{}
	if usr.Email != "" {
		filled = true
		match = append(match, bson.M{"email": usr.Email})
	}
	if usr.UserName != "" {
		filled = true
		match = append(match, bson.M{"username": usr.UserName})
	}
	if usr.ID != "" {
		// try find user by id or nickname
		match = append(match, bson.M{"id": usr.ID})
	}
	if !filled && usr.ID != "" {
		filled = true
		match = append(match, bson.M{"id": usr.ID})
	}

	if !filled {
		l.Error(zap.Error(ErrCannotFind), zap.Duration("duration", time.Since(start)), "incorrect data to search user")
		return nil, ErrCannotFind
	}

	u := &dto.UserDTO{}
	err := r.col.FindOne(ctx, bson.M{"$or": match}).Decode(u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "users record not found in database")
			return nil, ErrNotFound
		}
		l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "find error")
		return nil, ErrInternal
	}

	l.Info(zap.Duration("duration", time.Since(start)), "successful find user")
	userMapper := mapper.UserMapper{}
	return userMapper.ToModel(u), nil
}

func (r *DefaultUserRepository) Update(ctx context.Context, usr *domain.User) error {
	return nil
}

// Exist returns activity boolean of userid record from database
func (r *DefaultUserRepository) Exist(ctx context.Context, id string) (bool, error) {
	l := r.logger.Sugar().With("Exist")
	start := time.Now()
	if id == "" {
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "user id is empty")
		return false, ErrInvalidArgument
	}

	err := r.col.FindOne(ctx, bson.M{"id": id}).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Info(zap.Duration("duration", time.Since(start)), "user not found")
			return false, nil
		}
		l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "check user exist error")
		return false, ErrInternal
	}

	l.Info(zap.Duration("duration", time.Since(start)), "successful check user exist")
	return true, nil
}
