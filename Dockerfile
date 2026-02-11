# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /src

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/wallet-service ./cmd/wallet-service

# Runtime stage
FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/wallet-service /app/wallet-service
COPY config.env /app/config.env

EXPOSE 8080

CMD ["/app/wallet-service"]