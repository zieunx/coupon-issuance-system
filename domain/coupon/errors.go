package coupon

import "errors"

var ErrCouponIssuanceLimitExceeded = errors.New("쿠폰 발급 한도를 초과했습니다")
