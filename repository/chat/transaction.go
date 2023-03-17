package chat

import (
	"context"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.uber.org/zap"
	"time"
)

func (r *DefaultChatRepository) withTransactionChatCreate(ctx context.Context, ch *domain.Chat, chatsUsers []interface{}) error {
	l := r.logger.Sugar().With("withTransactionChatCreate")
	start := time.Now()
	if ch == nil {
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "chat is nil")
		return ErrEmpty
	}

	if ch == nil {
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "empty chat document")
		return ErrEmpty
	}

	if len(chatsUsers) == 0 {
		// The check is needed in case somehow it turned out that there were no participants, or the documents were not drawn up
		l.Error(zap.Error(ErrEmpty), zap.Duration("duration", time.Since(start)), "empty documents for the transaction")
		return ErrEmpty
	}

	opts := options.Session().SetDefaultReadConcern(readconcern.Local())
	sess, err := r.client.StartSession(opts)
	if err != nil {
		l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "failed start session")
		return err
	}

	defer sess.EndSession(ctx)

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Local()
	opt := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	err = mongo.WithSession(ctx, sess, func(sessionContext mongo.SessionContext) error {
		if err = sess.StartTransaction(opt); err != nil {
			l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "failed start transaction")
			return err
		}

		_, errInsertChat := r.col.InsertOne(sessionContext, ch)
		if errInsertChat != nil {
			l.Error(zap.Error(errInsertChat), zap.Duration("duration", time.Since(start)), "failed insert to chat collection")
			return errInsertChat
		}

		_, errInsertChatUser := r.colChatsUsers.InsertMany(sessionContext, chatsUsers)
		if errInsertChatUser != nil {
			l.Error(zap.Error(errInsertChatUser), zap.Duration("duration", time.Since(start)), "failed insert to chats_users collection")
			return errInsertChatUser
		}

		if err = sess.CommitTransaction(sessionContext); err != nil {
			l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "failed commit transaction")
			return err
		}

		return nil
	})

	if err != nil {
		l.Error(zap.Error(err), zap.Duration("duration", time.Since(start)), "failed start transaction")
		if abortErr := sess.AbortTransaction(ctx); abortErr != nil {
			l.Error(zap.Error(abortErr), zap.Duration("duration", time.Since(start)), "failed abort transaction")
			return abortErr
		}

		return err
	}

	return nil
}
