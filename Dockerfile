# --- Build stage ---
FROM golang:1.24-alpine AS builder

ARG TARGET
WORKDIR /app

RUN apk add --no-cache curl tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 빌드 대상 디렉토리로 이동
WORKDIR /app/${TARGET}
RUN go build -o main .

# --- Runtime stage ---
FROM alpine:latest
RUN apk add --no-cache curl bash tzdata

WORKDIR /app
ARG TARGET
COPY --from=builder /app/${TARGET}/main .

CMD ["./main"]
