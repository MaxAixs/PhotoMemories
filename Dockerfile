FROM golang:1.24-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -o mybot  ./cmd/main.go

FROM alpine:latest AS final
WORKDIR /root/
COPY --from=builder /app/mybot ./
COPY --from=builder /app/internal/config /root/internal/config
COPY .env .env
CMD ["/root/mybot"]