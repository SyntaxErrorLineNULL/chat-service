package user

import (
	"context"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
)

type IUserRepository interface {
	Create(ctx context.Context, usr *domain.User) error
	Find(ctx context.Context, usr *domain.User) (*domain.User, error)
	Update(ctx context.Context, usr *domain.User) error
}
