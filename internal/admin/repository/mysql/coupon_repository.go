package mysql

import (
	"context"
	"database/sql"
	"log"

	"coupon-issuance-system/domain/coupon"
)

type CouponRepositoryMySQL struct {
	db *sql.DB
}

func NewCouponRepositoryMySQL(db *sql.DB) coupon.CouponRepository {
	return &CouponRepositoryMySQL{db: db}
}

func (r *CouponRepositoryMySQL) GetCouponsByCampaignID(
	ctx context.Context,
	campaignID string,
) ([]*coupon.Coupon, error) {
	query := `SELECT ` +
		`id, code, issued_at, created_at, updated_at, user_id, campaign_id ` +
		`FROM coupon ` +
		`WHERE campaign_id = ?`

	log.Printf("Executing query: %s with campaignID: %s", query, campaignID)

	rows, err := r.db.QueryContext(ctx, query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coupons []*coupon.Coupon
	count := 0

	for rows.Next() {
		var coupon coupon.Coupon
		if err := rows.Scan(
			&coupon.ID,
			&coupon.Code,
			&coupon.IssuedAt,
			&coupon.CreatedAt,
			&coupon.UpdatedAt,
			&coupon.UserID,
			&coupon.CampaignID,
		); err != nil {
			return nil, err
		}
		coupons = append(coupons, &coupon)
		count++
	}

	log.Printf("Fetched %d coupons for campaignID: %s", count, campaignID)

	return coupons, nil
}

func (r *CouponRepositoryMySQL) CreateCoupon(
	ctx context.Context,
	coupon *coupon.Coupon,
) (*string, error) {
	// TODO: 구현하지 않음
	log.Println("CreateCoupon method is not implemented in CouponRepositoryMySQL")
	return nil, nil
}
