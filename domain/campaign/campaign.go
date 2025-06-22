package campaign

import "time"

type Campaign struct {
	ID                string    // 캠페인 고유 UUID
	Name              string    // 캠페인 이름
	CouponIssueLimit  int32     // 발급 가능한 총 수량
	IssuanceStartTime time.Time // 발급 시작 일시
	CreatedAt         time.Time // 생성 일시
	UpdatedAt         time.Time // 업데이트 일시
}
