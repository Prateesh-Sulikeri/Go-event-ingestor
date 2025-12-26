package limiter

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed token_bucket.lua
var luaScript string


type Limiter struct {
	client  *redis.Client
	max     int
	refill  int
}

func NewLimiter(client *redis.Client, maxTokens, refillRate int) *Limiter {
	return &Limiter{
		client: client,
		max:    maxTokens,
		refill: refillRate,
	}
}

// Allow returns true if the request is allowed, false if rate-limited.
func (l *Limiter) Allow(ctx context.Context, clientID string) (bool, error) {
	now := time.Now().Unix()
	
	bucketKey := "bucket:" + clientID
	timestampKey := "bucket_ts:" + clientID
	
	
	res, err := l.client.Eval(
    	ctx,
    	luaScript,
    	[]string{bucketKey, timestampKey},
    	l.max,
    	l.refill,
    	now,
	).Int()

	if err != nil {
    	return false, err
	}

	return res == 1, nil
}
