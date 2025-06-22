package client

import (
	"context"

	adminv1 "coupon-issuance-system/gen/admin/v1"
	adminv1connect "coupon-issuance-system/gen/admin/v1/adminv1connect"

	"connectrpc.com/connect"

	"coupon-issuance-system/domain/campaign"
)

type CampaignClient interface {
	GetCampaignByID(ctx context.Context, id string) (*campaign.Campaign, error)
}

type campaignClient struct {
	grpcClient adminv1connect.CampaignServiceClient
}

func NewCampaignClient(grpcClient adminv1connect.CampaignServiceClient) CampaignClient {
	return &campaignClient{
		grpcClient: grpcClient,
	}
}

func (c *campaignClient) GetCampaignByID(ctx context.Context, id string) (*campaign.Campaign, error) {
	req := connect.NewRequest(&adminv1.GetSimpleCampaignRequest{CampaignId: id})
	res, err := c.grpcClient.GetSimpleCampaign(ctx, req)
	if err != nil {
		return nil, err
	}
	return &campaign.Campaign{
		ID:                res.Msg.CampaignId,
		Name:              res.Msg.Name,
		CouponIssueLimit:  res.Msg.CouponIssueLimit,
		IssuanceStartTime: res.Msg.IssuanceStartTime.AsTime(),
	}, nil
}
