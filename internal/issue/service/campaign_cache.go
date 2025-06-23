package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"coupon-issuance-system/domain/campaign"
	"coupon-issuance-system/internal/issue/client"

	"github.com/redis/go-redis/v9"
)

// CampaignCache는 캠페인 정보를 가져오는 인터페이스입니다.
type CampaignCache interface {
	GetCampaign(ctx context.Context, campaignID string) (*campaign.Campaign, error)
}

type campaignCache struct {
	campaignClient client.CampaignClient
	redisClient    *redis.Client
	localCache     sync.Map
}

type localCachedCampaign struct {
	campaign *campaign.Campaign
	expires  time.Time
}

// NewCampaignCache는 새로운 CampaignCache를 생성합니다.
func NewCampaignCache(
	campaignClient client.CampaignClient,
	redisClient *redis.Client,
) CampaignCache {
	return &campaignCache{
		campaignClient: campaignClient,
		redisClient:    redisClient,
	}
}

// GetCampaign은 캐시(로컬, Redis)를 확인하고, 없는 경우 CampaignClient를 통해 캠페인 정보를 가져옵니다.
func (c *campaignCache) GetCampaign(ctx context.Context, campaignID string) (*campaign.Campaign, error) {
	// 1. 로컬 캐시 확인
	if val, ok := c.localCache.Load(campaignID); ok {
		if cached, ok := val.(localCachedCampaign); ok && time.Now().Before(cached.expires) {
			return cached.campaign, nil
		}
		c.localCache.Delete(campaignID)
	}
	redisKey := fmt.Sprintf("campaign:%s", campaignID)
	raw, err := c.redisClient.Get(ctx, redisKey).Result()

	if err == nil {
		var camp campaign.Campaign
		if err := json.Unmarshal([]byte(raw), &camp); err == nil {
			c.localCache.Store(campaignID, localCachedCampaign{
				campaign: &camp,
				expires:  time.Now().Add(1 * time.Second),
			})
			return &camp, nil
		}
	} else if err == redis.Nil {
		// Redis에 해당 키 없음 → 로깅
		log.Printf("Redis cache miss for campaign ID: %s", campaignID)
	} else {
		// Redis 조회 에러 → 에러 리턴
		return nil, fmt.Errorf("failed to get campaign from Redis: %w", err)
	}

	// 3. CampaignClient 호출
	camp, err := c.campaignClient.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	// Redis 캐시에 저장
	campBytes, err := json.Marshal(camp)
	if err == nil {
		c.redisClient.Set(ctx, redisKey, campBytes, 10*time.Minute)
	}

	// 로컬 캐시에 저장
	c.localCache.Store(campaignID, localCachedCampaign{
		campaign: camp,
		expires:  time.Now().Add(1 * time.Second),
	})

	return camp, nil
}
