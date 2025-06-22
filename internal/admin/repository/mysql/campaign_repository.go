package mysql

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"coupon-issuance-system/domain/campaign"

	"github.com/google/uuid"
)

type CampaignRepositoryMySQL struct {
	db *sql.DB
}

func NewCampaignRepositoryMySQL(db *sql.DB) campaign.CampaignRepository {
	return &CampaignRepositoryMySQL{db: db}
}

// CreateCampaign 새로운 캠페인 생성
func (r *CampaignRepositoryMySQL) CreateCampaign(
	ctx context.Context,
	campaign *campaign.Campaign,
) (*string, error) {
	campaignID := uuid.New().String()

	query := `INSERT INTO campaign ` +
		`(id, name, coupon_issue_limit, issuance_start_time) ` +
		`VALUES (?, ?, ?, ?)`

	log.Printf("Executing query: %s with values: %s, %s, %d, %s",
		query,
		campaignID,
		campaign.Name,
		campaign.CouponIssueLimit,
		campaign.IssuanceStartTime,
	)

	_, err := r.db.ExecContext(ctx, query,
		campaignID,
		campaign.Name,
		campaign.CouponIssueLimit,
		campaign.IssuanceStartTime,
	)
	if err != nil {
		return nil, err
	}

	return &campaignID, nil
}

func (r *CampaignRepositoryMySQL) GetCampaignByID(
	ctx context.Context,
	id string,
) (*campaign.Campaign, error) {
	query := `SELECT ` +
		`id, name, coupon_issue_limit, issuance_start_time, created_at, updated_at ` +
		`FROM campaign ` +
		`WHERE id = ?`

	log.Printf("Executing query: %s with id: %s", query, id)

	row := r.db.QueryRowContext(ctx, query, id)
	var foundCampaign campaign.Campaign
	err := row.Scan(
		&foundCampaign.ID,
		&foundCampaign.Name,
		&foundCampaign.CouponIssueLimit,
		&foundCampaign.IssuanceStartTime,
		&foundCampaign.CreatedAt,
		&foundCampaign.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, campaign.ErrCampaignNotFound
		}
		return nil, err
	}
	return &foundCampaign, nil
}
