package mysql

import (
	"context"
	"database/sql"
	"log"

	"coupon-issuance-system/domain/coupon"

	"github.com/google/uuid"
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
		`id, code, issued_at, created_at, updated_at, campaign_id ` +
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
			&coupon.CampaignID,
		); err != nil {
			return nil, err
		}
		coupons = append(coupons, &coupon)
		count++
	}

	return coupons, nil
}

func (r *CouponRepositoryMySQL) CreateCoupon(
	ctx context.Context,
	coupon *coupon.Coupon,
) (*string, error) {
	couponID := uuid.New().String()

	query := `INSERT INTO coupon (id, code, issued_at, user_id, campaign_id) ` +
		`VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		couponID,
		coupon.Code,
		coupon.IssuedAt,
		coupon.UserID,
		coupon.CampaignID,
	)
	if err != nil {
		return nil, err
	}

	return &couponID, nil
}
