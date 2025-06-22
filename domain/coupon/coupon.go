package coupon

import "time"

type Coupon struct {
	ID         string    // 쿠폰 고유 UUID
	Code       string    // 쿠폰 코드
	IssuedAt   time.Time // 쿠폰 발급 일시
	CreatedAt  time.Time // 쿠폰 생성 일시
	UpdatedAt  time.Time // 쿠폰 업데이트 일시
	CampaignID string    // 쿠폰이 속한 캠페인 ID
	UserID     string    // 쿠폰을 발급받은 사용자 ID
}
