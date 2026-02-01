# Stage 1: Build Frontend
FROM oven/bun:latest AS frontend-builder
WORKDIR /app
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile
COPY . .
RUN bun run build

# Stage 2: Build Backend
FROM golang:1.25.5-alpine AS backend-builder
WORKDIR /app
# 复制 go.mod 和 go.sum 优先下载依赖，利用镜像层缓存
COPY go.mod go.sum ./
RUN go mod download
# 复制所有源代码
COPY . .
# 从前端构建阶段复制 dist 目录，用于 go:embed
COPY --from=frontend-builder /app/dist ./dist
# 禁用 CGO，添加 -tags=go_json，并使用 ldflags 减小体积
RUN CGO_ENABLED=0 GOOS=linux go build -tags=go_json -ldflags="-s -w" -o main .

# Stage 3: Runtime
# 使用 scratch 或 alpine。如果需要时区或证书，alpine 更方便
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
ENV GIN_MODE=release

WORKDIR /app
# 仅复制二进制文件
COPY --from=backend-builder /app/main .

EXPOSE 8080
CMD ["./main"]
