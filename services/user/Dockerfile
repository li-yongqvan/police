FROM docker.m.daocloud.io/library/golang:1.23-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ./cmd/main.go

FROM docker.m.daocloud.io/library/alpine:3.19
RUN apk --no-cache add ca-certificates wget
WORKDIR /root/
COPY --from=builder /user-service .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8001
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8001/health || exit 1
CMD ["./user-service"]
