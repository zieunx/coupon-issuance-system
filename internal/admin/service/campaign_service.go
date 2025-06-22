package service

import (
	"context"
	"errors"

	"coupon-issuance-system/domain/campaign"
	"coupon-issuance-system/domain/coupon"

	"connectrpc.com/connect"
)

type CampaignService interface {
	CreateCampaign(ctx context.Context, req *CreateCampaignRequest) (*CampaignMsg, error)
	GetCampaign(ctx context.Context, req *GetCampaignRequest) (*CampaignMsg, error)
	GetSimpleCampaign(ctx context.Context, req *GetCampaignRequest) (*CampaignMsg, error)
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
