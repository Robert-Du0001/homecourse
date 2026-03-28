# --- 前端构建 ---
FROM node:24-alpine AS node-builder
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

RUN apk add --no-cache tzdata

WORKDIR /www

COPY --from=node-builder /node-build/public/ ./public/

COPY --from=builder /build/main /build/entrypoint.sh ./
COPY --from=builder /build/storage/ ./storage/
COPY --from=builder /build/.env.example ./.env

RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["/bin/sh", "/www/entrypoint.sh"]
