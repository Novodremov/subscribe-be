SERVICE_INTERFACE=internal/service/iface.go
SERVICE_MOCK_DIR=internal/service/mocks
SERVICE_MOCK_DEST=$(SERVICE_MOCK_DIR)/mock_service.go

REPO_INTERFACE=internal/repo/iface.go
REPO_MOCK_DIR=internal/repo/mocks
REPO_MOCK_DEST=$(REPO_MOCK_DIR)/mock_repo.go

MOCK_PKG=mocks

install-mockgen:
	go install github.com/golang/mock/mockgen@latest

mocks: service-mocks repo-mocks

service-mocks:
	mkdir -p $(SERVICE_MOCK_DIR)
	mockgen -source=$(SERVICE_INTERFACE) -destination=$(SERVICE_MOCK_DEST) -package=$(MOCK_PKG)

repo-mocks:
	mkdir -p $(REPO_MOCK_DIR)
	mockgen -source=$(REPO_INTERFACE) -destination=$(REPO_MOCK_DEST) -package=$(MOCK_PKG)

test:
	CGO_ENABLED=1 go test -count=10 -race ./...

cover:
	go test -short -count=1 \
		-coverpkg=$$(go list ./... \
		| grep -v '/config' \
		| grep -v '/repo' \
		| grep -v '/mocks' \
		| grep -v '/docs' \
		| grep -v '/app' \
		| grep -v '/db' \
		| grep -v '/logging' \
		| tr '\n' ',') \
    	-coverprofile=coverage.out ./...

	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

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
