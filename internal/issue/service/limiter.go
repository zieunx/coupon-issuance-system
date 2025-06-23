package service

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Limiter는 발급 수량을 제한하는 인터페이스
type Limiter interface {
	// Allow는 발급이 허용되는지 확인
	Allow(ctx context.Context, campaignID string, limit int) (bool, error)
	// Rollback은 발급 실패 시 카운터를 롤백
	Rollback(ctx context.Context, campaignID string) error
}

type redisLimiter struct {
	redisClient *redis.Client
}

func NewRedisLimiter(redisClient *redis.Client) Limiter {
	return &redisLimiter{redisClient: redisClient}
}

// Allow는 Redis를 사용하여 발급 수량을 확인하고, 한도를 초과하지 않은 경우 true를 반환합니다.
func (l *redisLimiter) Allow(ctx context.Context, campaignID string, limit int) (bool, error) {
	key := fmt.Sprintf("campaign:%s:issued_count", campaignID)
	count, err := l.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis INCR failed: %w", err)
	}

	if count > int64(limit) {
		l.redisClient.Decr(ctx, key)
		return false, nil
	}

	return true, nil
}

// Rollback은 Redis 카운터를 감소시킵니다.
func (l *redisLimiter) Rollback(ctx context.Context, campaignID string) error {
	key := fmt.Sprintf("campaign:%s:issued_count", campaignID)
	return l.redisClient.Decr(ctx, key).Err()
}
