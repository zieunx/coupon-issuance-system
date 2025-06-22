package service

import (
	"context"

	"coupon-issuance-system/domain/campaign"
	"coupon-issuance-system/domain/coupon"
)

type CampaignService interface {
	CreateCampaign(ctx context.Context, req *CreateCampaignRequest) (*CampaignMsg, error)
	GetCampaign(ctx context.Context, req *GetCampaignRequest) (*CampaignMsg, error)
}

type campaignService struct {
	CampaignRepository campaign.CampaignRepository
	CouponRepository   coupon.CouponRepository
}

func NewCampaignService(
	campaignRepository campaign.CampaignRepository,
	couponRepository coupon.CouponRepository,
) CampaignService {
	return &campaignService{
		CampaignRepository: campaignRepository,
		CouponRepository:   couponRepository,
	}
}

func (s *campaignService) CreateCampaign(ctx context.Context, req *CreateCampaignRequest) (*CampaignMsg, error) {
	newCampaign := &campaign.Campaign{
		Name:              req.Name,
		CouponIssueLimit:  req.CouponIssueLimit,
		IssuanceStartTime: req.IssuanceStartTime,
	}

	_, err := s.CampaignRepository.CreateCampaign(ctx, newCampaign)
	if err != nil {
		return nil, err
	}

	return &CampaignMsg{
		ID:                newCampaign.ID,
		CouponIssueLimit:  newCampaign.CouponIssueLimit,
		IssuanceStartTime: newCampaign.IssuanceStartTime,
		CreatedAt:         newCampaign.CreatedAt,
		UpdatedAt:         newCampaign.UpdatedAt,
	}, nil
}

func (s *campaignService) GetCampaign(ctx context.Context, req *GetCampaignRequest) (*CampaignMsg, error) {
	// 캠페인 조회
	campaign, err := s.CampaignRepository.GetCampaignByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// 쿠폰들 조회
	coupons, err := s.CouponRepository.GetCouponsByCampaignID(ctx, campaign.ID)
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
		ID:                campaign.ID,
		Name:              campaign.Name,
		CouponIssueLimit:  campaign.CouponIssueLimit,
		IssuanceStartTime: campaign.IssuanceStartTime,
		CreatedAt:         campaign.CreatedAt,
		UpdatedAt:         campaign.UpdatedAt,
		Coupons:           couponMsgs,
	}, nil
}
