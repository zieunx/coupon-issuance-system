package handler

import (
	"context"
	issuev1 "coupon-issuance-system/gen/issue/v1"
	issuev1connect "coupon-issuance-system/gen/issue/v1/issuev1connect"

	connect "connectrpc.com/connect"

	"coupon-issuance-system/internal/issue/service"
)

type IssueHandler struct {
	service service.IssueService
}

func NewIssueHandler(svc service.IssueService) issuev1connect.IssueServiceHandler {
	return &IssueHandler{service: svc}
}

func (h *IssueHandler) IssueCoupon(
	ctx context.Context,
	req *connect.Request[issuev1.IssueCouponRequest],
) (*connect.Response[issuev1.IssueCouponResponse], error) {
	in := req.Msg

	// 내부 요청으로 변환
	internalReq := &service.IssueCouponRequest{
		CampaignID: in.GetCampaignId(),
		UserID:     in.GetUserId(),
	}

	// 서비스 호출
	couponId, err := h.service.IssueCoupon(ctx, internalReq)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 응답 변환
	resp := &issuev1.IssueCouponResponse{
		CouponId: *couponId,
	}

	return connect.NewResponse(resp), nil
}
