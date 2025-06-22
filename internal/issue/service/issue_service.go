package service

import (
	"context"
	"time"

	"coupon-issuance-system/domain/campaign"
	"coupon-issuance-system/domain/coupon"
	"coupon-issuance-system/util"
)

type IssueService interface {
	IssueCoupon(ctx context.Context, req *IssueCouponRequest) (*string, error)
}

type issueService struct {
	CampaignRepository campaign.CampaignRepository
	CouponRepository   coupon.CouponRepository
}

func NewIssueService(
	campaignRepository campaign.CampaignRepository,
	couponRepository coupon.CouponRepository,
) IssueService {
	return &issueService{
		CampaignRepository: campaignRepository,
		CouponRepository:   couponRepository,
	}
}

func (s *issueService) IssueCoupon(ctx context.Context, req *IssueCouponRequest) (*string, error) {
	// 캠페인 조회
	_, err := s.CampaignRepository.GetCampaignByID(ctx, req.CampaignID)
	if err != nil {
		return nil, err
	}

	// 쿠폰 발급 로직 구현
	newCoupon := &coupon.Coupon{
		Code:       util.GenerateCouponCode(),
		IssuedAt:   time.Now(),
		CampaignID: req.CampaignID,
		UserID:     req.UserID,
	}

	// 쿠폰 저장
	couponId, err := s.CouponRepository.CreateCoupon(ctx, newCoupon)
	if err != nil {
		return nil, err
	}

	return couponId, nil
}
