package user

import (
	"context"
	"errors"
	"fmt"
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

// ErrUserEmptyID returns when userid is empty
var ErrUserEmptyID = errors.New("empty user id")

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

func (srv *Service) methodLogger(ctx context.Context, name string) *zap.Logger {
	return srv.logger.With(
		zap.String("layer: {user-service} ", fmt.Sprintf("method: [%s]", name)),
	)
}

// Create a new user and save it to the database
// This function creates a new user by checking if the input data is valid, generating a unique ID, and saving the user to the database.
// It logs errors and returns error messages if anything goes wrong.
func (srv *Service) Create(ctx context.Context, input *domain.User) error {
	// Get the method logger with additional context information.
	l := srv.methodLogger(ctx, "Create")
	l.Debug("creating new chat")

	// Check if the input data is empty
	if input == nil {
		l.Error("empty request data", zap.Error(ErrEmpty))
		return ErrEmpty
	}

	// check Email validation use net mail pkg
	if !input.ValidEmail() {
		l.Error("failed check email validation", zap.Error(ErrIncorrectEmail))
		return ErrIncorrectEmail
	}

	// Check if the user already exists
	_, err := srv.Find(ctx, input)
	if err == nil {
		l.Error("user already exist", zap.Error(ErrEmpty))
		return ErrUserAlreadyExist
	}

	// Generate a unique ID for the new user
	id := uuid.New().String()
	input.ID = id

	err = srv.db.User.Create(ctx, input)
	if err != nil {
		l.Error("failed create user", zap.Error(err))
		return err
	}

	l.Info("successful create user")
	// Return nil if everything succeeded
	return nil
}

// Find searches for a user in the database using the provided context and input.
// Returns the user if found, otherwise returns an error.
func (srv *Service) Find(ctx context.Context, input *domain.User) (*domain.User, error) {
	// Get the method logger with additional context information.
	l := srv.methodLogger(ctx, "Find")
	l.Debug("find user record")

	// Check if the input is empty
	if input == nil {
		l.Error("empty request data", zap.Error(ErrEmpty))
		return nil, ErrEmpty
	}

	// Find the user in the database
	u, err := srv.db.User.Find(ctx, input)
	// Check if the user was not found
	if errors.Is(err, user.ErrNotFound) || errors.Is(err, user.ErrCannotFind) {
		l.Error("user not found", zap.Error(ErrUserNotFound))
		return input, ErrUserNotFound
	}
	// Check if there was an error while finding the user
	if err != nil {
		l.Error("failed find user", zap.Error(err))
		return input, err
	}

	l.Info("successful find user record")
	return u, err
}

// Update updates an existing user record in the database.
// It takes a context and a user domain object as input,
// and returns an error if any operation fails.
func (srv *Service) Update(ctx context.Context, input *domain.User) error {
	// Get the method logger with additional context information.
	l := srv.methodLogger(ctx, "Update")
	l.Debug("update user record")

	// Check if user data is empty
	if input == nil {
		l.Error("empty user data", zap.Error(ErrEmpty))
		return ErrEmpty
	}

	// Call the Update method of the User repository to update the user record
	err := srv.db.User.Update(ctx, input)
	if err != nil {
		l.Error("failed update user", zap.Error(err))
		return err
	}

	// Log success and return nil
	l.Info("successfully update user data")
	return nil
}

// Exist checks if a user with the specified ID exists in the database.
// It returns a boolean indicating whether the user exists or not, along with an error, if any.
func (srv *Service) Exist(ctx context.Context, userID string) (bool, error) {
	// Get the method logger with additional context information.
	l := srv.methodLogger(ctx, "Exist")
	l.Debug("check user existence")

	// Check if the user ID is empty.
	if userID == "" {
		l.Error("empty user data", zap.Error(ErrEmpty))
		return false, ErrUserEmptyID
	}

	// Call the database to check if the user exists.
	exist, err := srv.db.User.Exist(ctx, userID)
	if err != nil {
		l.Error("failed check user exist", zap.Error(err))
		return false, err
	}

	// Log success and return the result.
	l.Info("successfully checked user exist")
	return exist, nil
}
