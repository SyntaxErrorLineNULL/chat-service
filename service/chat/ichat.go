package chat

import (
	"context"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
)

type IChat interface {
	Create(ctx context.Context, input *domain.Chat) error
}
