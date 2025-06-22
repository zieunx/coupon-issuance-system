package campaign

import (
	"context"
)

type CampaignRepository interface {
	CreateCampaign(ctx context.Context, campaign *Campaign) (*string, error)
	GetCampaignByID(ctx context.Context, id string) (*Campaign, error)
}
