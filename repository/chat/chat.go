package chat

import (
	"context"
	"errors"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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

// Find for chat by id or by id and participant or by id and chat owner.
func (r *DefaultChatRepository) Find(ctx context.Context, ch *domain.Chat) (*domain.Chat, error) {
	l := r.logger.Sugar().With("Create")
	start := time.Now()
	if ch == nil {
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "chat is nil")
		return nil, ErrEmpty
	}

	chat := &domain.Chat{}

	filled := false
	match := bson.A{}
	if len(ch.Participants) != 0 {
		filled = true
		// Array must contain all passed values. Operator $all solves this problem.
		// https://www.mongodb.com/docs/manual/reference/operator/query/all/#mongodb-query-op.-all
		match = append(match, bson.M{"participants": bson.M{"$all": ch.Participants}})
	}
	if ch.OwnerID != "" {
		filled = true
		match = append(match, bson.M{"owner_id": ch.OwnerID})
	}
	if ch.ID != "" {
		filled = true
		match = append(match, bson.M{"id": ch.ID})
	}

	if !filled {
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "incorrect data to find chat")
		return nil, ErrCannotFind
	}

	err := r.col.FindOne(ctx, bson.M{"$and": match}).Decode(chat)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "chat not found in database")
			return nil, ErrNotFound
		}
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "find error")
		return nil, err
	}

	return chat, nil
}

// FindOwnedChats returns all chats where the user is the owner
func (r *DefaultChatRepository) FindOwnedChats(ctx context.Context, ownerID string) ([]*domain.Chat, error) {
	l := r.logger.Sugar().With("FindActiveChatsByOwner")
	start := time.Now()
	if ownerID == "" {
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "empty owner id")
		return nil, ErrEmpty
	}

	cursor, err := r.col.Find(ctx, bson.M{"owner_id": ownerID}, nil)
	if err != nil {
		l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "failed find chat with owner id")
		return nil, err
	}

	var ch []*domain.Chat
	if err = cursor.All(ctx, &ch); err != nil {
		l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "failed find all documents")
		return nil, err
	}

	if err = cursor.Close(ctx); err != nil {
		l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "failed close cursor")
		return nil, err
	}

	return ch, nil
}

// FindPersonalChatBetweenUsers chat document
func (r *DefaultChatRepository) FindPersonalChatBetweenUsers(ctx context.Context, ownerID, participantID string) (*domain.Chat, error) {
	l := r.logger.Sugar().With("FindPersonalChatBetweenUsers")
	start := time.Now()
	if ownerID == "" || participantID == "" {
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "empty owner and participant id")
		return nil, ErrEmpty
	}

	chat := &domain.Chat{}
	filter := bson.D{
		{Key: "owner_id", Value: ownerID},
		{Key: "participants", Value: participantID},
		{Key: "type", Value: domain.ChatTypePersonal},
	}
	err := r.col.FindOne(ctx, filter).Decode(chat)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			l.Error(zap.Error(ErrNotFound), zap.Duration("duration", time.Since(start)), "not found in database")
			return nil, ErrNotFound
		}
		l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "find chat error")
		return nil, err
	}

	return chat, nil
}
