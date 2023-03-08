package chat

import (
	"context"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
)

type IChatRepository interface {
	Create(ctx context.Context, ch *domain.Chat) error
	Find(ctx context.Context, ch *domain.Chat) (*domain.Chat, error)
	FindOwnedChats(ctx context.Context, ownerID string) ([]*domain.Chat, error)
}
