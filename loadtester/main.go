package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	issuev1 "coupon-issuance-system/gen/issue/v1"
	"coupon-issuance-system/gen/issue/v1/issuev1connect"

	"connectrpc.com/connect"
)

func main() {
	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		IdleConnTimeout:     30, // optional
	}
	httpClient := &http.Client{
		Transport: transport,
	}
	client := issuev1connect.NewIssueServiceClient(
		httpClient,
		"http://issue-server:8081",
	)

	total := 2000
	success := 0
	failed := 0
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req := &issuev1.IssueCouponRequest{
				CampaignId: "abc123",
				UserId:     fmt.Sprintf("user-%d", i),
			}
			res, err := client.IssueCoupon(context.Background(), connect.NewRequest(req))
			mutex.Lock()
			defer mutex.Unlock()
			if err != nil || res.Msg.CouponId == "" {
				failed++
				return
			}
			success++
		}(i)
	}

	wg.Wait()
	fmt.Printf("✔️ 성공: %d\n❌ 실패: %d\n", success, failed)
}
