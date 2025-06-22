package coupon

import (
	"context"
)

type CouponRepository interface {
	GetCouponsByCampaignID(ctx context.Context, campaignID string) ([]*Coupon, error)
	CreateCoupon(ctx context.Context, coupon *Coupon) (*string, error)
}
