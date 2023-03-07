package chat

import (
	"context"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type DefaultChatRepository struct {
	client        *mongo.Client
	col           *mongo.Collection
	colChatsUsers *mongo.Collection
	logger        *zap.Logger
}

func NewDefaultChatRepository(client *mongo.Client, logger *zap.Logger) *DefaultChatRepository {
	return &DefaultChatRepository{
		client:        client,
		col:           client.Database("chat-service").Collection("chat"),
		colChatsUsers: client.Database("chat-service").Collection("chats_users"),
		logger:        logger,
	}
}

// Create  inserts new chat document
func (r *DefaultChatRepository) Create(ctx context.Context, ch *domain.Chat) error {
	l := r.logger.Sugar().With("Create")
	if ch == nil {
		return ErrEmpty
	}
	start := time.Now()
	if ch == nil {
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "chat is nil")
		return ErrEmpty
	}

	t := time.Now().UnixMilli()
	ch.CreateDate = t

	tt := time.Now().AddDate(10, 0, 0).UnixMilli()
	participants := make([]interface{}, 0)
	for _, e := range ch.Participants {
		chatsUsers := &domain.ChatsUsers{
			ID:             uuid.New().String(),
			ChatID:         ch.ID,
			UserID:         e,
			AddedAt:        t,
			StartMessageID: t,
			EndMessageID:   tt,
			MaxReadDate:    t,
		}
		participants = append(participants, chatsUsers)
	}

	// Start transaction
	if err := r.withTransactionChatCreate(ctx, ch, participants); err != nil {
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "failed insert chat with transaction")
		return err
	}

	return nil
}
