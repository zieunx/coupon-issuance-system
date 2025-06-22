package main

import (
	"fmt"
	"log"
	"net/http"

	adminv1connect "coupon-issuance-system/gen/admin/v1/adminv1connect"
	"coupon-issuance-system/internal/admin/handler"
	"coupon-issuance-system/internal/admin/repository/mysql"
	"coupon-issuance-system/internal/admin/service"
	"coupon-issuance-system/internal/config"
	"coupon-issuance-system/internal/config/database"
	"coupon-issuance-system/internal/interceptor"

	connect "connectrpc.com/connect"
)

func main() {
	// 설정 로드
	config.LoadConfig()
	cfg := config.GetConfig()

	// DB 연결
	db, err := database.ConnectDB(&cfg.MySQL)
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}

	// 테이블 생성
	if err := database.EnsureTables(db); err != nil {
		log.Fatalf("테이블 생성 실패: %v", err)
	}

	// 의존성 주입
	campaignRepository := mysql.NewCampaignRepositoryMySQL(db)
	couponRepository := mysql.NewCouponRepositoryMySQL(db)
	svc := service.NewCampaignService(campaignRepository, couponRepository)
	server := handler.NewCampaignHandler(svc)

	// ConnectRPC 핸들러 설정
	path, h := adminv1connect.NewCampaignServiceHandler(
		server,
		connect.WithInterceptors(interceptor.NewLoggingInterceptor()),
	)

	// HTTP 서버 설정
	mux := http.NewServeMux()
	mux.Handle(path, h)

	// 서버 시작
	addr := fmt.Sprintf(":%s", cfg.AdminServer.Port)
	log.Printf("쿠폰 관리 시스템 서버 시작: %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("서버 종료: %v", err)
	}
}
