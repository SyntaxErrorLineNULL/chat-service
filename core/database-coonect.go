package core

import (
	"context"
	"github.com/SyntaxErrorLineNULL/chat-service/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

func DatabaseConnect(ctx context.Context, cfg *config.Config, logger zap.Logger) (*mongo.Client, error) {
	log := logger.With().With(zap.Field{
		Key:    "core",
		String: "DatabaseConnect",
	})
	start := time.Now()
	log.Info("init connect to database", zap.Duration("duration", time.Since(start)))
	mongoOption := options.Client().ApplyURI(cfg.Database.URL).SetLoadBalanced(true)
	client, err := mongo.Connect(ctx, mongoOption)
	if err != nil {
		log.Error("failed connect to database", zap.Field{Interface: err})
		return nil, err
	}

	// ping database connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Error("failed ping database connection", zap.Field{Interface: err})
		return nil, err
	}

	return client, err
}
