package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"coupon-issuance-system/domain/coupon"
	"coupon-issuance-system/util"

	"connectrpc.com/connect"
)

type IssueService interface {
	IssueCoupon(ctx context.Context, req *IssueCouponRequest) (*string, error)
}

type issueService struct {
	campaignCache CampaignCache
	limiter       Limiter
	repo          coupon.CouponRepository
}

func NewIssueService(
	campaignCache CampaignCache,
	limiter Limiter,
	repo coupon.CouponRepository,
) IssueService {
	return &issueService{
		campaignCache: campaignCache,
		limiter:       limiter,
		repo:          repo,
	}
}

func (s *issueService) IssueCoupon(ctx context.Context, req *IssueCouponRequest) (*string, error) {
	// 캠페인 정보 조회
	campaign, err := s.campaignCache.GetCampaign(ctx, req.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	now := time.Now()
	if now.Before(campaign.IssuanceStartTime) {
		return nil, connect.NewError(
			connect.CodeFailedPrecondition,
			fmt.Errorf("쿠폰 발일시급 가능 는 '%s'입니다.", campaign.IssuanceStartTime.Format(time.RFC3339)),
		)
	}

	// 캠페인 쿠폰 발급 제한 확인
	allowed, err := s.limiter.Allow(ctx, campaign.ID, int(campaign.CouponIssueLimit))
	if err != nil {
		return nil, fmt.Errorf("failed to check issuance limit: %w", err)
	}
	if !allowed {
		return nil, coupon.ErrCouponIssuanceLimitExceeded
	}

	// 새로운 쿠폰 생성
	newCoupon := &coupon.Coupon{
		Code:       util.GenerateCouponCode(),
		IssuedAt:   time.Now(),
		CampaignID: campaign.ID,
		UserID:     req.UserID,
	}

	couponId, err := s.repo.CreateCoupon(ctx, newCoupon)
	if err != nil {
		log.Default().Println("Failed to create coupon:", err)
		if rErr := s.limiter.Rollback(ctx, campaign.ID); rErr != nil {
			log.Default().Println("Failed to rollback counter:", rErr)
		}
		return nil, err
	}

	return couponId, nil
}
