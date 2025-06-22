package util

import (
	"math/rand"
	"sync"
	"time"
)

// 한글 28자 + 숫자 10자
var charset = []rune("가나다라마바사아자차카타파하거너더러머버서어저처커터퍼허0123456789")

const codeLength = 7

var charsetLen = int64(len(charset))

// 고루틴 안전한 rand.Seed 사용을 위한 sync
var (
	mu     sync.Mutex
	seeded = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// 쿠폰 코드 생성 함수
func GenerateCouponCode() string {
	mu.Lock()
	defer mu.Unlock()

	var code []rune
	for i := 0; i < codeLength; i++ {
		idx := seeded.Int63() % charsetLen
		code = append(code, charset[idx])
	}

	return string(code)
}
