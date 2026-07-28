package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var likeRateLimitScript = redis.NewScript(`
local count = redis.call('INCR', KEYS[1])
local window = tonumber(ARGV[1])

if count == 1 then
    redis.call('PEXPIRE', KEYS[1], window)
end

local ttl = redis.call('PTTL', KEYS[1])
if ttl < 0 then
    redis.call('PEXPIRE', KEYS[1], window)
    ttl = window
end

return {count, ttl}
`)

type LikeLimiter struct {
	client redis.Scripter
	limit  int64
	window time.Duration
}

func NewLikeLimiter(
	client redis.Scripter,
	limit int64,
	window time.Duration,
) *LikeLimiter {
	return &LikeLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (l *LikeLimiter) Allow(
	ctx context.Context,
	userID uint64,
) (bool, time.Duration, error) {
	key := fmt.Sprintf("rate_limit:like:user:%d", userID)

	result, err := likeRateLimitScript.Run(
		ctx,
		l.client,
		[]string{key},
		l.window.Milliseconds(),
	).Int64Slice()
	if err != nil {
		return false, 0, fmt.Errorf("执行点赞限流脚本失败: %w", err)
	}
	if len(result) != 2 {
		return false, 0, fmt.Errorf("点赞限流脚本返回值无效")
	}

	count := result[0]
	retryAfter := time.Duration(result[1]) * time.Millisecond

	return count <= l.limit, retryAfter, nil
}
