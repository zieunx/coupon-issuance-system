package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

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

	// Redis 연결
	rdb, err := database.ConnectRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("Redis 연결 실패: %v", err)
	}

	// 테이블 생성
	if err := database.EnsureTables(db); err != nil {
		log.Fatalf("테이블 생성 실패: %v", err)
	}

	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     30 * time.Second,
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   3 * time.Second, // 요청당 타임아웃 (적절히 설정)
	}

	adminGRPCClient := adminv1connect.NewCampaignServiceClient(
		httpClient,
		fmt.Sprintf("http://%s:%s", cfg.AdminServer.Host, cfg.AdminServer.Port),
	)

	campaignClient := client.NewCampaignClient(adminGRPCClient)
	couponRepository := mysql.NewCouponRepositoryMySQL(db)

	campaignCache := service.NewCampaignCache(campaignClient, rdb)
	limiter := service.NewRedisLimiter(rdb)

	svc := service.NewIssueService(campaignCache, limiter, couponRepository)
	issueHandler := handler.NewIssueHandler(svc)

	// ConnectRPC 핸들러 설정
	path, h := issuev1connect.NewIssueServiceHandler(
		issueHandler,
		connect.WithInterceptors(interceptor.NewLoggingInterceptor()),
	)

	// HTTP mux 설정
	mux := http.NewServeMux()
	mux.Handle(path, h)

	// ✅ 서버 리스너 명시적으로 생성 (backlog 제어 가능)
	addr := fmt.Sprintf(":%s", cfg.IssueServer.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("서버 Listen 실패: %v", err)
	}

	server := &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("쿠폰 발급 시스템 서버 시작: %s", addr)
	if err := server.Serve(ln); err != nil {
		log.Fatalf("서버 종료: %v", err)
	}
}
