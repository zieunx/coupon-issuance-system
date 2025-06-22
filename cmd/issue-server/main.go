package main

import (
	"fmt"
	"log"
	"net/http"

	"coupon-issuance-system/gen/admin/v1/adminv1connect"
	issuev1connect "coupon-issuance-system/gen/issue/v1/issuev1connect"
	"coupon-issuance-system/internal/config"
	"coupon-issuance-system/internal/config/database"
	"coupon-issuance-system/internal/interceptor"
	"coupon-issuance-system/internal/issue/client"
	"coupon-issuance-system/internal/issue/handler"
	"coupon-issuance-system/internal/issue/repository/mysql"
	"coupon-issuance-system/internal/issue/service"

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
	// gRPC 클라이언트 생성
	adminGRPCClient := adminv1connect.NewCampaignServiceClient(
		http.DefaultClient,
		"http://localhost:8081", // 실제 Admin 서버 주소
	)
	campaignClient := client.NewCampaignClient(adminGRPCClient)
	couponRepository := mysql.NewCouponRepositoryMySQL(db)
	svc := service.NewIssueService(campaignClient, couponRepository)
	server := handler.NewIssueHandler(svc)

	// ConnectRPC 핸들러 설정
	path, h := issuev1connect.NewIssueServiceHandler(
		server,
		connect.WithInterceptors(interceptor.NewLoggingInterceptor()),
	)

	// HTTP 서버 설정
	mux := http.NewServeMux()
	mux.Handle(path, h)

	// 서버 시작
	addr := fmt.Sprintf(":%s", cfg.IssueServer.Port)
	log.Printf("쿠폰 발급 시스템 서버 시작: %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("서버 종료: %v", err)
	}
}
