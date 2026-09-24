FROM golang:1.27 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o bot \
    ./cmd/bot

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

COPY certs/russian_trusted_root_ca.crt /usr/local/share/ca-certificates/russian_trusted_root_ca.crt
RUN update-ca-certificates

WORKDIR /app

COPY --from=builder --chmod=0755 /app/bot .

CMD ["./bot"]
