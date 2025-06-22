package handler

import (
	"context"
	"time"

	adminv1 "coupon-issuance-system/gen/admin/v1"
	adminv1connect "coupon-issuance-system/gen/admin/v1/adminv1connect"

	connect "connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"coupon-issuance-system/internal/admin/service"
)

type CampaignHandler struct {
	service service.CampaignService
}

func NewCampaignHandler(svc service.CampaignService) adminv1connect.CampaignServiceHandler {
	return &CampaignHandler{service: svc}
}

func (s *CampaignHandler) CreateCampaign(
	ctx context.Context,
	req *connect.Request[adminv1.CreateCampaignRequest],
) (*connect.Response[adminv1.CreateCampaignResponse], error) {
	in := req.Msg

	// proto 요청을 내부 요청으로 변환
	internalReq := &service.CreateCampaignRequest{
		Name:              in.GetName(),
		CouponIssueLimit:  in.GetCouponIssueLimit(),
		IssuanceStartTime: in.GetIssuanceStartTime().AsTime(),
	}

	// 순수한 비즈니스 로직 호출 (proto 타입 없음)
	campaignMsg, err := s.service.CreateCampaign(ctx, internalReq)
	if err != nil {
		// 에러 타입에 따라 적절한 ConnectRPC 에러로 변환
		if _, ok := err.(*time.ParseError); ok {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 내부 응답을 proto 응답으로 변환
	protoResp := &adminv1.CreateCampaignResponse{
		CampaignId: campaignMsg.ID,
	}
	return connect.NewResponse(protoResp), nil
}
func (s *CampaignHandler) GetCampaign(
	ctx context.Context,
	req *connect.Request[adminv1.GetCampaignRequest],
) (*connect.Response[adminv1.GetCampaignResponse], error) {
	in := req.Msg

	// proto 요청을 내부 요청으로 변환
	internalReq := &service.GetCampaignRequest{
		ID: in.GetCampaignId(),
	}

	// 비즈니스 로직 호출
	campaignMsg, err := s.service.GetCampaign(ctx, internalReq)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// proto 응답 생성
	protoResp := &adminv1.GetCampaignResponse{
		CampaignId:        campaignMsg.ID,
		CouponIssueLimit:  int32(campaignMsg.CouponIssueLimit),
		IssuanceStartTime: timestamppb.New(campaignMsg.IssuanceStartTime),
		Coupons:           make([]*adminv1.CouponResponse, len(campaignMsg.Coupons)),
	}

	// 쿠폰 목록 변환
	for i, c := range campaignMsg.Coupons {
		protoResp.Coupons[i] = &adminv1.CouponResponse{
			Id:         c.ID,
			Code:       c.Code,
			IssuedAt:   timestamppb.New(c.IssuedAt),
			CreatedAt:  timestamppb.New(c.CreatedAt),
			UpdatedAt:  timestamppb.New(c.UpdatedAt),
			CampaignId: c.CampaignID,
		}
	}

	return connect.NewResponse(protoResp), nil
}
