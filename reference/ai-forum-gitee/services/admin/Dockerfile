FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /admin-service ./cmd/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates wget
WORKDIR /root/
COPY --from=builder /admin-service .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8003
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8003/health || exit 1
CMD ["./admin-service"]
