package service

import "time"

type CreateCampaignRequest struct {
	Name              string    // 캠페인 이름
	CouponIssueLimit  int32     // 발급 가능한 총 수량
	IssuanceStartTime time.Time // 발급 시작 일시
}

type GetCampaignRequest struct {
	ID string
}

type CampaignMsg struct {
	ID                string       // 캠페인 고유 ID
	Name              string       // 캠페인 이름
	CouponIssueLimit  int32        // 발급 가능한 총 수량
	IssuanceStartTime time.Time    // 발급 시작 일시
	CreatedAt         time.Time    // 생성 일시
	UpdatedAt         time.Time    // 업데이트 일시
	Coupons           []*CouponMsg // 쿠폰 목록
}

type CouponMsg struct {
	ID         string    // 쿠폰 고유 UUID
	Code       string    // 쿠폰 코드
	IssuedAt   time.Time // 쿠폰 발급 일시
	CreatedAt  time.Time // 쿠폰 생성 일시
	UpdatedAt  time.Time // 쿠폰 업데이트 일시
	CampaignID string    // 쿠폰이 속한 캠페인 ID
}
