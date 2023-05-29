package user

import (
	"context"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
)

type UserService interface {
	Create(ctx context.Context, input *domain.User) error
	Find(ctx context.Context, input *domain.User) (*domain.User, error)
	Update(ctx context.Context, input *domain.User) error
	Exist(ctx context.Context, userID string) (bool, error)
}
