package interceptor

import (
	"context"
	"log"

	connect "connectrpc.com/connect"
)

// NewLoggingInterceptor는 ConnectRPC 요청/응답 로깅 인터셉터를 생성
func NewLoggingInterceptor() connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			log.Printf("[ConnectRPC] %s %s %s", req.Spec().Procedure, req.Peer().Addr, req.Any())
			resp, err := next(ctx, req)
			if err != nil {
				log.Printf("[ConnectRPC] ERROR: %v", err)
			} else {
				log.Printf("[ConnectRPC] OK: %s", req.Spec().Procedure)
			}
			return resp, err
		})
	})
}
