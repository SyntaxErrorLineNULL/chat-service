package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SyntaxErrorLineNULL/chat-service/domain"
	"github.com/SyntaxErrorLineNULL/chat-service/repository/dto"
	"github.com/SyntaxErrorLineNULL/chat-service/repository/mapper"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// ErrInvalidArgument returns when request was with invalid data
var ErrInvalidArgument = errors.New("invalid argument")

type DefaultUserRepository struct {
	client *mongo.Client
	col    *mongo.Collection
	logger *zap.Logger
}

func NewDefaultUserRepository(client *mongo.Client, logger *zap.Logger) *DefaultUserRepository {
	l := logger.With(zap.String("repository", "confirmation"))
	return &DefaultUserRepository{
		client: client,
		col:    client.Database("chat-service").Collection("user"),
		logger: l,
	}
}

func (r *DefaultUserRepository) methodLogger(ctx context.Context, name string) *zap.Logger {
	return r.logger.With(
		zap.String("layer: {user-repository} ", fmt.Sprintf("method: [%s]", name)),
	)
}

// Create inserts a new user into the database.
func (r *DefaultUserRepository) Create(ctx context.Context, usr *domain.User) error {
	// Get the method logger with additional context information.
	l := r.methodLogger(ctx, "Create")

	// Log that we're starting the creation of a new user.
	l.Debug("creating new user")

	// Start a timer to measure how long it takes to create the user.
	start := time.Now()

	// Check if the provided user object is nil. If it is, log an error and return an error.
	if usr == nil {
		l.Error("user is nil", zap.Error(ErrInvalidArgument), zap.Duration("duration", time.Since(start)))
		return ErrInvalidArgument
	}

	// If the provided user doesn't have an ID, generate a new UUID.
	if usr.ID == "" {
		usr.ID = uuid.New().String()
	}

	// Map the user domain object to a DTO object and insert it into the database.
	u := &mapper.UserMapper{}
	_, err := r.col.InsertOne(ctx, u.ToDTO(usr))
	if err != nil {
		// If there was an error inserting the user, log an error and return an internal server error.
		l.Error("insert error", zap.Error(err), zap.Duration("duration", time.Since(start)))
		return ErrInternal
	}

	// If we made it this far, the user was successfully created. Log an info message and return nil.
	l.Info("successfully created user", zap.Duration("duration", time.Since(start)))
	return nil
}

// Find searches for a user in the database based on provided criteria
// and returns a domain.User object representing the found user, or an error.
func (r *DefaultUserRepository) Find(ctx context.Context, usr *domain.User) (*domain.User, error) {
	// Get the method logger with additional context information.
	l := r.methodLogger(ctx, "Find")
	l.Debug("find user record in database")

	start := time.Now()

	// Check if user is nil
	if usr == nil {
		l.Error("user is nil", zap.Error(ErrInvalidArgument), zap.Duration("duration", time.Since(start)))
		return nil, ErrInvalidArgument
	}

	filled := false
	match := bson.A{}
	// Build search criteria based on the provided user object.
	// If email, username or id are provided, search using them.
	// If none of them are provided, return an error.
	if usr.Email != "" {
		filled = true
		match = append(match, bson.M{"email": usr.Email})
	}
	if usr.UserName != "" {
		filled = true
		match = append(match, bson.M{"username": usr.UserName})
	}
	if !filled && usr.ID != "" {
		filled = true
		match = append(match, bson.M{"id": usr.ID})
	}

	if !filled {
		l.Error("incorrect data to search user", zap.Error(ErrCannotFind), zap.Duration("duration", time.Since(start)))
		return nil, ErrCannotFind
	}

	u := &dto.UserDTO{}
	err := r.col.FindOne(ctx, bson.M{"$or": match}).Decode(u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If no user was found, return ErrNotFound error
			l.Error("users record not found in database", zap.Error(err), zap.Duration("duration", time.Since(start)))
			return nil, ErrNotFound
		}
		// If there was an error with the database, return ErrInternal error
		l.Error("find error", zap.Error(err), zap.Duration("duration", time.Since(start)))
		return nil, ErrInternal
	}

	// If user was found, return it as a domain.User object
	l.Info("successful find user", zap.Duration("duration", time.Since(start)))
	userMapper := mapper.UserMapper{}
	return userMapper.ToModel(u), nil
}

// Update user document
func (r *DefaultUserRepository) Update(ctx context.Context, usr *domain.User) error {
	// Get the method logger with additional context information.
	l := r.methodLogger(ctx, "Update")
	l.Debug("update user record")
	start := time.Now()
	if usr == nil {
		l.Error("user is nil", zap.Error(ErrInvalidArgument), zap.Duration("duration", time.Since(start)))
		return ErrInvalidArgument
	}

	u := &mapper.UserMapper{}
	data := u.ToDTO(usr)
	_, err := r.col.UpdateOne(
		ctx,
		bson.M{"id": data.ID},
		bson.M{"$set": data},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Error("users record not found in database", zap.Error(err), zap.Duration("duration", time.Since(start)))
			return ErrNotFound
		}
		l.Error("update user data error", zap.Error(err), zap.Duration("duration", time.Since(start)))
		return ErrInternal
	}

	l.Info("successful update user data", zap.Duration("duration", time.Since(start)))
	return nil
}

// Exist returns activity boolean of userid record from database
func (r *DefaultUserRepository) Exist(ctx context.Context, id string) (bool, error) {
	// Get the method logger with additional context information.
	l := r.methodLogger(ctx, "Exist")
	l.Debug("check exist record in database")
	// Start the timer to measure the method execution time
	start := time.Now()

	// Return an error if the user ID is empty
	if id == "" {
		l.Error("user id is empty", zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)))
		return false, ErrInvalidArgument
	}

	// Check if the user with the given ID exists in the database
	err := r.col.FindOne(ctx, bson.M{"id": id}).Err()
	if err != nil {
		// Return false if the user is not found
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Info("user not found", zap.Duration("duration", time.Since(start)))
			return false, nil
		}
		// Return an error if an unexpected error occurs while checking for user existence
		l.Error("check user exist error", zap.Error(err), zap.Duration("duration", time.Since(start)))
		return false, ErrInternal
	}

	// Log and return true if the user exists in the database
	l.Info("successful check user exist", zap.Duration("duration", time.Since(start)))
	return true, nil
}

// ExistUserName returns boolean check the existence of the record in the database
func (r *DefaultUserRepository) ExistUserName(ctx context.Context, username string) (bool, error) {
	// Get the method logger with additional context information.
	l := r.methodLogger(ctx, "v")
	// Debug log to indicate the start of the method
	l.Debug("check exist username in database")
	start := time.Now()
	if username == "" {
		// If username is empty, log an error and return ErrInvalidArgument
		l.Error("username is empty", zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)))
		return false, ErrInvalidArgument
	}

	// Check if a user document exists with the provided username
	err := r.col.FindOne(ctx, bson.M{"username": username}).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If no documents match the query, log an info message and return false
			l.Info("username not exist", zap.Duration("duration", time.Since(start)))
			return false, nil
		}
		// If an error other than ErrNoDocuments occurs, log an error and return ErrInternal
		l.Error("check username exist error", zap.Error(err), zap.Duration("duration", time.Since(start)))
		return false, ErrInternal
	}

	// If a matching document is found, log an info message and return true
	l.Info("successful check username exist", zap.Duration("duration", time.Since(start)))
	return true, nil
}
