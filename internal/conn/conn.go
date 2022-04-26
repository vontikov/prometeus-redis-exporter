package conn

import (
	"context"
	"sync"

	redis "github.com/go-redis/redis/v8"

	"github.com/vontikov/prom-redis/internal/collector"
	"github.com/vontikov/prom-redis/internal/logging"
)

type RedisWrapper struct {
	ctx         context.Context
	opts        *redis.Options
	infoSection string
	logger      logging.Logger

	mu     sync.Mutex // protects following fields
	client *redis.Client
}

func New(ctx context.Context, opts *redis.Options, infoSection string) *RedisWrapper {
	return &RedisWrapper{
		ctx:         ctx,
		opts:        opts,
		infoSection: infoSection,
		logger:      logging.NewLogger("redis"),
	}
}

func (r *RedisWrapper) Importer() collector.Importer {
	return func() *string {
		r.mu.Lock()
		defer r.mu.Unlock()

		if r.client == nil {
			r.client = redis.NewClient(r.opts)
		}
		res, err := r.client.Info(r.ctx, r.infoSection).Result()
		if err != nil {
			r.logger.Error("error", "message", err)
			r.client = nil
			return nil
		}
		return &res
	}
}
