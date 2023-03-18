package user

import (
	"context"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
)

// UserRepository database interface
type UserRepository interface {
	Create(ctx context.Context, usr *domain.User) error
	Find(ctx context.Context, usr *domain.User) (*domain.User, error)
	Update(ctx context.Context, usr *domain.User) error
	Exist(ctx context.Context, id string) (bool, error)
	ExistUserName(ctx context.Context, username string) (bool, error)
}
