package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"coupon-issuance-system/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

// ConnectDB는 DB 연결을 담당
func ConnectDB(config *config.MySQLConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(80)                 // 동시에 열 수 있는 최대 커넥션
	db.SetMaxIdleConns(40)                 // 대기 커넥션 수 (보통 max의 50%)
	db.SetConnMaxLifetime(5 * time.Minute) // 오래된 커넥션 정리

	if err := db.PingContext(context.Background()); err != nil {
		return nil, err
	}

	log.Printf("DB 연결 성공: %s", config.DSN)

	return db, nil
}

// ConnectRedis는 Redis 연결을 담당
func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	log.Printf("Redis 연결 성공: %s", cfg.Address)
	return rdb, nil
}

// EnsureTables는 필요한 테이블을 생성
func EnsureTables(db *sql.DB) error {
	CreateCampaignQuery := `
	CREATE TABLE IF NOT EXISTS campaign (
	  id VARCHAR(255) PRIMARY KEY,
	  name VARCHAR(255) NOT NULL,
	  coupon_issue_limit INT NOT NULL,
	  issuance_start_time TIMESTAMP NOT NULL,
	  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(CreateCampaignQuery); err != nil {
		return fmt.Errorf("failed to create campaign table: %w", err)
	}

	CreateCouponQuery := `
	CREATE TABLE IF NOT EXISTS coupon (
		id VARCHAR(255) PRIMARY KEY,
		code VARCHAR(255) NOT NULL UNIQUE,
		issued_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		user_id VARCHAR(255),
		campaign_id VARCHAR(255),
		FOREIGN KEY (campaign_id) REFERENCES campaign(id)
	);`

	if _, err := db.Exec(CreateCouponQuery); err != nil {
		return fmt.Errorf("failed to create coupon table: %w", err)
	}

	return nil
}
