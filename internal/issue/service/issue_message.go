package service

type IssueCouponRequest struct {
	CampaignID string // 캠페인 ID
	UserID     string // 사용자 ID
}

type IssueCouponResponse struct {
	ID   string // 쿠폰 고유 UUID
	Code string // 쿠폰 코드
}
