FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/subscribe-be ./cmd/...

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/subscribe-be /app/subscribe-be
COPY --from=builder /app/db/migrations /app/db/migrations
COPY .env /app/.env

EXPOSE 8080

CMD ["/app/subscribe-be"]
