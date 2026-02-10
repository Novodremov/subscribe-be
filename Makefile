env:
	cp .env.example .env
.PHONY: env

tidy:
	go mod tidy
.PHONY: tidy

sqlc:
	sqlc generate
.PHONY: sqlc

swagger:
	swag init -g internal/delivery/http/router.go -o ./assets/docs
