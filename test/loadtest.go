package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"connectrpc.com/connect"

	issuev1 "coupon-issuance-system/gen/issue/v1"
	issuev1connect "coupon-issuance-system/gen/issue/v1/issuev1connect"
)

func main() {
	var (
		serverAddr  string
		campaignID  string
		total       int
		concurrency int
	)
	flag.StringVar(&serverAddr, "server", "http://localhost:8082", "RPC server address")
	flag.StringVar(&campaignID, "campaign", "", "Campaign ID to test against")
	flag.IntVar(&total, "total", 500, "Total number of requests to send")
	flag.IntVar(&concurrency, "concurrency", 500, "Max concurrency level")
	flag.Parse()

	if campaignID == "" {
		fmt.Println("❗ 캠페인 ID는 필수입니다. 예: -campaign=cmp_abc")
		os.Exit(1)
	}

	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
	}
	httpClient := &http.Client{
		Transport: transport,
	}
	client := issuev1connect.NewIssueServiceClient(
		httpClient,
		serverAddr,
	)

	start := make(chan struct{})
	var wg sync.WaitGroup
	var mu sync.Mutex

	success := 0
	fail := 0
	couponCodes := make(map[string]bool)
	semaphore := make(chan struct{}, concurrency)

	wg.Add(total)
	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()
			<-start

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			req := connect.NewRequest(&issuev1.IssueCouponRequest{
				CampaignId: campaignID,
			})

			resp, err := client.IssueCoupon(context.Background(), req)

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				fail++
				log.Printf("❌ 요청 실패: %v", err)
			} else {
				code := resp.Msg.GetCouponId()
				if couponCodes[code] {
					log.Printf("⚠️ 중복 쿠폰 코드: %s", code)
				}
				couponCodes[code] = true
				success++
			}
		}()
	}

	fmt.Printf("🔥 %d개의 요청 준비 완료 (동시성 %d), 1초 후 시작!\n", total, concurrency)
	time.Sleep(1 * time.Second)
	startTime := time.Now()
	close(start)

	wg.Wait()
	elapsed := time.Since(startTime)

	fmt.Println("✅ 부하 테스트 완료")
	fmt.Printf("총 요청 수: %d\n", total)
	fmt.Printf("성공: %d, 실패: %d, 중복 코드 없음: %t\n", success, fail, len(couponCodes) == success)
	fmt.Printf("실제 발급된 고유 쿠폰 수: %d\n", len(couponCodes))
	fmt.Printf("총 소요 시간: %v\n", elapsed)
}
