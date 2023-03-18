package user

import (
	"context"
	"errors"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
	"github.com/SyntaxErrorLineNULL/chat-service/repository"
	"github.com/SyntaxErrorLineNULL/chat-service/repository/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ErrEmpty returns when request was with empty data
var ErrEmpty = errors.New("empty request")

// ErrUserNotFound returns when cannot find user in database
var ErrUserNotFound = errors.New("user not found")

// ErrUserAlreadyExist returns when user already exists in the database.
// Username and email must be unique
var ErrUserAlreadyExist = errors.New("user already exist with the same credentials")

// ErrIncorrectEmail returns case the email is incorrect
var ErrIncorrectEmail = errors.New("email is incorrect")

// Settings user service setting provider
type Settings struct {
	DB     *repository.Repository
	Logger *zap.Logger
}

// Service user ?? mb i rename this)
type Service struct {
	db     *repository.Repository
	logger *zap.Logger
}

// New returns new instance of User service
func New(settings *Settings) *Service {
	l := settings.Logger
	l.Sugar().With("New user service")
	return &Service{
		db:     settings.DB,
		logger: l,
	}
}

// Create a new user and save it to the database
// This function creates a new user by checking if the input data is valid, generating a unique ID, and saving the user to the database.
// It logs errors and returns error messages if anything goes wrong.
func (srv *Service) Create(ctx context.Context, input *domain.User) error {
	// Log the start of the method call
	l := srv.logger.Sugar().With("Create")
	l.Debug("creating new chat")

	// Check if the input data is empty
	if input == nil {
		l.Error(zap.Error(ErrEmpty), "empty request data")
		return ErrEmpty
	}

	// check Email validation use net mail pkg
	if !input.ValidEmail() {
		l.Error(zap.Error(ErrIncorrectEmail), "failed check email validation")
		return ErrIncorrectEmail
	}

	// Check if the user already exists
	_, err := srv.Find(ctx, input)
	if err == nil {
		l.Error(zap.Error(ErrEmpty), "user already exist")
		return ErrUserAlreadyExist
	}

	// Generate a unique ID for the new user
	id := uuid.New().String()
	input.ID = id

	err = srv.db.User.Create(ctx, input)
	if err != nil {
		l.Error(zap.Error(err), "failed create user")
		return err
	}

	l.Info("successful create user")
	// Return nil if everything succeeded
	return nil
}

// Find searches for a user in the database using the provided context and input.
// Returns the user if found, otherwise returns an error.
func (srv *Service) Find(ctx context.Context, input *domain.User) (*domain.User, error) {
	l := srv.logger.Sugar().With("Find")
	l.Debug("find user record")

	// Check if the input is empty
	if input == nil {
		l.Error(zap.Error(ErrEmpty), "empty request data")
		return nil, ErrEmpty
	}

	// Find the user in the database
	u, err := srv.db.User.Find(ctx, input)
	// Check if the user was not found
	if errors.Is(err, user.ErrNotFound) || errors.Is(err, user.ErrCannotFind) {
		l.Error(zap.Error(ErrEmpty), "user not found")
		return input, ErrUserNotFound
	}
	// Check if there was an error while finding the user
	if err != nil {
		l.Error(zap.Error(ErrEmpty), "failed find user")
		return input, err
	}

	return u, err
}
