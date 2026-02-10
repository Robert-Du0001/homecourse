# --- 前端构建 ---
FROM node:20-alpine AS node-builder
# 安装 pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /node-build
COPY . .
RUN pnpm config set registry https://registry.npmmirror.com && \
  pnpm install && \
  pnpm run build

# --- 后端构建 ---
FROM golang:alpine AS builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0

WORKDIR /build
COPY . .

RUN go mod tidy
RUN go build --ldflags "-s -w -extldflags -static" -o main .

FROM alpine:latest

WORKDIR /www

COPY --from=node-builder /node-build/public/ /www/public/

COPY --from=builder /build/main /www/
COPY --from=builder /build/storage/ /www/storage/
COPY --from=builder /build/.env.example /www/.env

RUN /www/main artisan key:generate && \
    /www/main artisan jwt:secret && \
    /www/main artisan migrate

ENTRYPOINT ["/www/main"]
