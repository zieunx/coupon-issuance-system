# coupon-issuance-system


## 🚀 빠른 시작

1. 의존 서비스 실행 (MySQL, Redis)
```bash
docker-compose up -d
```

2. Admin 서버 실행

```bash
go run ./cmd/admin-server
```

3. Issue 서버 실행

```bash
go run ./cmd/issue-server
```



## 🛠️ API 

### 1) 캠페인 생성

```bash
curl -X POST http://localhost:8081/admin.v1.CampaignService/CreateCampaign \
  -H "Content-Type: application/json" \
  -H "Connect-Protocol-Version: 1" \
  -d '{
    "name": "캠페인1",
    "couponIssueLimit": 1000,
    "issuanceStartTime": "2025-06-23T20:57:00+09:00"
  }'
```

### 2) 캠페인 조회

```bash
curl -X POST http://localhost:8081/admin.v1.CampaignService/GetCampaign \
  -H "Content-Type: application/json" \
  -H "Connect-Protocol-Version: 1" \
  -d '{
    "campaignId": "생성된_campaignId_값"
  }'

```


### 3) 쿠폰 발급

```bash
curl -X POST http://localhost:8082/issue.v1.IssueService/IssueCoupon \
  -H "Content-Type: application/json" \
  -H "Connect-Protocol-Version: 1" \
  -d '{
    "campaignId": "생성된_campaignId_값",
    "userId": "test"
  }'

```

## 설계


![스크린샷 2025-06-23 오후 9 30 59](https://github.com/user-attachments/assets/b78fac80-87aa-4852-b6a3-c46344c50d66)

### 1. 캠페인 정보 캐싱 (with Redis)

- **관리 주체**: Admin Server
- **시나리오**
    1. Admin Server에서 캠페인 등록·수정·삭제 API 호출
    2. 내부 로직 처리 후 Redis에 캠페인 데이터 전체를 `HMSET`으로 저장
    3. Issue Server는 오직 Redis 조회(`HGETALL`)만 사용
- **특징**
    - TTL 없음 → 수동으로 갱신된 데이터만 존재
    - 캐시 미스 시(미구현)
        - Admin Server로 페일백 요청 후 Redis 재저장

```mermaid
flowchart LR
  A[Admin Server] -->|등록/수정/삭제| B[Redis Cache]
  C[Issue Server] -->|조회| B
```


### 2. 쿠폰 코드 생성기

- **위치**: `./util/code_generator.go`
- **목표**: 한글+숫자 조합 10자리 안전 랜덤 문자열 생성

| 항목          | 값                                    |
| ------------- | ------------------------------------- |
| 문자 집합     | 한글 36자 + 숫자 10자 = 총 46자       |
| 코드 길이     | 10자리                                |
| 조합 수       | 46¹⁰ ≈ 3.58×10¹⁶                      |
| 중복 확률(1M) | ≈ 0.0014% (사실상 중복 없음)           |


### 3. Redis 원자 연산 (Atomic INCR/DECR)

- 목적: 캠페인별로 발급 가능한 쿠폰 수량을 초과하지 않도록, 모든 서버 인스턴스에서 동시성 문제 없이 안전하게 제어합니다.
- 구현 방식:
  - 쿠폰 발급 시, Redis의 INCR 명령어로 해당 캠페인의 발급 카운터(campaign:{id}:issued_count)를 1 증가
  - 증가된 값이 캠페인에 설정된 최대 발급 수량을 초과하면, 즉시 DECR로 롤백하여 초과 발급을 방지
  - 실제 쿠폰 DB 저장에 실패한 경우에도, DECR로 카운터를 복구(롤백)
  - 이 과정은 모든 서버에서 원자적으로 처리되므로, 수평 확장 환경에서도 데이터 정합성이 보장
- 코드 참고:
  - internal/issue/service/limiter.go
  - internal/issue/service/issue_service.go

