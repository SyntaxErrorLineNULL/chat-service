package core

import (
	"context"
	"github.com/SyntaxErrorLineNULL/chat-service/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

func RedisConnect(ctx context.Context, cfg *config.Config, logger zap.Logger) error {
	log := logger.With().With(zap.Field{
		Key:    "core",
		String: "RedisConnect",
	})
	start := time.Now()
	log.Info("Init create redis cluster connect", zap.Duration("duration", time.Since(start)))
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:      cfg.Redis.Address,
		Password:   cfg.Redis.Password,
		ClientName: cfg.Redis.ClientName,
	})

	err := rdb.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		log.Info("Ping each redis cluster shard", zap.Duration("duration", time.Since(start)))
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		log.Error("failed ping shard", zap.Field{Interface: err})
		return err
	}

	return nil
}
