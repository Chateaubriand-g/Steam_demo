FROM go:1.24-alpine AS builder
WORKDIR /app
RUN echo "https://mirrors.aliyun.com/alpine/v3.18/main/" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.18/community/" >> /etc/apk/repositories && \
    # 更新软件源索引（刷新包列表）
    apk update && \
    # 清理缓存（避免使用损坏的旧包）
    rm -rf /var/cache/apk/* && \
    apk add --no-cache git

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -p 1 -buildvcs=false -o steam-backend ./api

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/steam-backend /app/
#COPY .env .env
EXPOSE 8080
CMD ["./steam-backend"]
