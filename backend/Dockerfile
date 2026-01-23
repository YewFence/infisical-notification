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

# 默认环境变量
ENV APP_ENV=prod
ENV TODO_DB_PATH=/app/data/todos.db

COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
