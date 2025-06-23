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
			// log.Printf("[ConnectRPC] %s %s %s", req.Spec().Procedure, req.Peer().Addr, req.Any())

			resp, err := next(ctx, req)

			// 네트워크/시스템 레벨 에러
			if err != nil {
				log.Printf("[ConnectRPC] ERROR: %s %s -> %v", req.Spec().Procedure, req.Peer().Addr, err)
				return resp, err
			}

			// Connect 에러는 실제로는 error가 아닌 응답 내부에 담긴 상태 코드 (단언 필요)
			if r, ok := resp.(interface{ Error() error }); ok {
				if r.Error() != nil {
					log.Printf("[ConnectRPC] APP ERROR: %s %s -> %v", req.Spec().Procedure, req.Peer().Addr, r.Error())
				}
			}

			// log.Printf("[ConnectRPC] OK: %s", req.Spec().Procedure)
			return resp, nil
		})
	})
}
