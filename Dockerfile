# Stage 1: Build Frontend
FROM oven/bun:latest AS frontend-builder
WORKDIR /app/web
COPY web/package.json web/bun.lock ./
RUN bun install --frozen-lockfile
COPY web/ ./
RUN bun run build
# Stage 2: Build Backend
FROM golang:1.25.5-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -tags=go_json -ldflags="-s -w" -o main .
# Stage 3: Runtime
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata mailcap
ENV TZ=Asia/Shanghai
ENV GIN_MODE=release
WORKDIR /app
COPY --from=backend-builder /app/main .
EXPOSE 8080
CMD ["./main"]