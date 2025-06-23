package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"coupon-issuance-system/domain/campaign"
	"coupon-issuance-system/domain/coupon"

	"connectrpc.com/connect"
	"github.com/redis/go-redis/v9"
)

type CampaignService interface {
	CreateCampaign(ctx context.Context, req *CreateCampaignRequest) (*CampaignMsg, error)
	GetCampaign(ctx context.Context, req *GetCampaignRequest) (*CampaignMsg, error)
	GetSimpleCampaign(ctx context.Context, req *GetCampaignRequest) (*CampaignMsg, error)
}

type campaignService struct {
	CampaignRepository campaign.CampaignRepository
	CouponRepository   coupon.CouponRepository
	RedisClient        *redis.Client
}

func NewCampaignService(
	campaignRepository campaign.CampaignRepository,
	couponRepository coupon.CouponRepository,
	redis *redis.Client,
) CampaignService {
	return &campaignService{
		CampaignRepository: campaignRepository,
		CouponRepository:   couponRepository,
		RedisClient:        redis,
	}
}

func (s *campaignService) CreateCampaign(ctx context.Context, req *CreateCampaignRequest) (*CampaignMsg, error) {
	newCampaign := &campaign.Campaign{
		Name:              req.Name,
		CouponIssueLimit:  req.CouponIssueLimit,
		IssuanceStartTime: req.IssuanceStartTime,
	}

	campaignId, err := s.CampaignRepository.CreateCampaign(ctx, newCampaign)
	if err != nil {
		return nil, err
	}

	newCampaign.ID = *campaignId // Redis 저장을 위해 ID 설정

	// Redis 저장 (TTL 없이)
	cacheKey := fmt.Sprintf("campaign:%s", *campaignId)
	campaignJSON, err := json.Marshal(newCampaign)
	if err != nil {
		return nil, err
	}
	if err := s.RedisClient.Set(ctx, cacheKey, campaignJSON, 0).Err(); err != nil {
		return nil, fmt.Errorf("redis 저장 실패: %w", err)
	}

	return &CampaignMsg{
		ID:                *campaignId,
		CouponIssueLimit:  newCampaign.CouponIssueLimit,
		IssuanceStartTime: newCampaign.IssuanceStartTime,
		CreatedAt:         newCampaign.CreatedAt,
		UpdatedAt:         newCampaign.UpdatedAt,
	}, nil
}

func (s *campaignService) GetCampaign(ctx context.Context, req *GetCampaignRequest) (*CampaignMsg, error) {
	// 캠페인 조회
	foundCampaign, err := s.CampaignRepository.GetCampaignByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, campaign.ErrCampaignNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, campaign.ErrCampaignNotFound)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 쿠폰들 조회
	coupons, err := s.CouponRepository.GetCouponsByCampaignID(ctx, foundCampaign.ID)
	if err != nil {
		return nil, err
	}

	// 쿠폰 메시지 변환
	var couponMsgs []*CouponMsg
	for _, coupon := range coupons {
		couponMsgs = append(couponMsgs, &CouponMsg{
			ID:         coupon.ID,
			Code:       coupon.Code,
			IssuedAt:   coupon.IssuedAt,
			CreatedAt:  coupon.CreatedAt,
			UpdatedAt:  coupon.UpdatedAt,
			UserID:     coupon.UserID,
			CampaignID: coupon.CampaignID,
		})
	}

	// 캠페인 메시지 반환
	return &CampaignMsg{
		ID:                foundCampaign.ID,
		Name:              foundCampaign.Name,
		CouponIssueLimit:  foundCampaign.CouponIssueLimit,
		IssuanceStartTime: foundCampaign.IssuanceStartTime,
		CreatedAt:         foundCampaign.CreatedAt,
		UpdatedAt:         foundCampaign.UpdatedAt,
		Coupons:           couponMsgs,
	}, nil

}

func (s *campaignService) GetSimpleCampaign(ctx context.Context, req *GetCampaignRequest) (*CampaignMsg, error) {
	foundCampaign, err := s.CampaignRepository.GetCampaignByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, campaign.ErrCampaignNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, campaign.ErrCampaignNotFound)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &CampaignMsg{
		ID:                foundCampaign.ID,
		Name:              foundCampaign.Name,
		CouponIssueLimit:  foundCampaign.CouponIssueLimit,
		IssuanceStartTime: foundCampaign.IssuanceStartTime,
		CreatedAt:         foundCampaign.CreatedAt,
		UpdatedAt:         foundCampaign.UpdatedAt,
	}, nil
}
