# 阶段1: 构建
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# 阶段2: 运行
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/server .

# 默认生产环境，可通过 docker compose 环境变量覆盖
ENV APP_ENV=prod

EXPOSE 8080
CMD ["./server"]
