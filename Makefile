GO_VERSION=$(shell wget -qO- "https://golang.org/VERSION?m=text" | grep -o 'go[0-9.]*')
CURRENT_GO_VERSION=$(shell go version | awk '{print $$3}')
GO_LINT=$(shell which golangci-lint 2> /dev/null || echo '')
GO_LINT_URI=github.com/golangci/golangci-lint/cmd/golangci-lint@latest
GO_SEC=$(shell which gosec 2> /dev/null || echo '')
GO_SEC_URI=github.com/securego/gosec/v2/cmd/gosec@latest
GO_VULNCHECK=$(shell which govulncheck 2> /dev/null || echo '')
GO_VULNCHECK_URI=golang.org/x/vuln/cmd/govulncheck@latest
SQLC_VERSION := $(shell sqlc version 2>/dev/null || echo "none")
SQLC_URI=github.com/sqlc-dev/sqlc/cmd/sqlc@latest
SQLC_DESTINATION=internal/transactions/adapters/repository/postgres/db
MIGRATE_VERSION := $(shell migrate -version 2>/dev/null || echo "none")
MIGRATE_URI=github.com/golang-migrate/migrate/v4/cmd/migrate@latest
GO_VULNCHECK=$(shell which govulncheck 2> /dev/null || echo '')
MOCK_GEN=$(shell which mockgen 2> /dev/null || echo '')
MOCK_GEN_URI=go.uber.org/mock/mockgen@latest
MOCK_DESTINATION=internal/mocks

.PHONY: run
run:
	go run cmd/main.go -file=$(file) -email=$(email)

.PHONY: update_go
update_go:
	@echo "Latest Go version: $(GO_VERSION)"
	@echo "Local Go version: $(CURRENT_GO_VERSION)"
	@if [ "$(CURRENT_GO_VERSION)" != "$(GO_VERSION)" ]; then \
		echo "Updating Go from $(CURRENT_GO_VERSION) to $(GO_VERSION)"; \
		wget https://golang.org/dl/$(GO_VERSION).linux-amd64.tar.gz; \
		sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $(GO_VERSION).linux-amd64.tar.gz; \
		rm $(GO_VERSION).linux-amd64.tar.gz; \
	else \
		echo "Go is already up to date"; \
	fi

.PHONY: docker-compose
docker-up:
	docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down

.PHONY: sqlc
sqlc:
	@if [ "$(SQLC_VERSION)" = "none" ]; then \
		echo "Installing sqlc..."; \
		go install $(SQLC_URI); \
	else \
		echo "sqlc is already installed"; \
	fi
	@rm -rf $(SQLC_DESTINATION)
	sqlc generate

.PHONY: migrate-new
migrate-new:
	@if [ "$(MIGRATE_VERSION)" = "none" ]; then \
		echo "Installing golang-migrate..."; \
		go install $(MIGRATE_URI); \
	else \
		echo "golang-migrate is already installed"; \
	fi
	migrate create -ext sql -dir db-scripts/migrations -seq $(NAME)

.PHONY: migrate-up
migrate-up:
	migrate -path db-scripts/migrations -database "postgres://root@localhost:5432/gostori?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	migrate -path db-scripts/migrations -database "postgres://root@localhost:5432/gostori?sslmode=disable" down

.PHONY: migrate-force
migrate-force:
	migrate -path db-scripts/migrations -database "postgres://root@localhost:5432/gostori?sslmode=disable" force $(VERSION)

.PHONY: lint
lint:
	$(if $(GO_LINT), ,go install $(GO_LINT_URI))
	golangci-lint run -v

.PHONY: sec
sec:
	$(if $(GO_SEC), ,cd /tmp && go install $(GO_SEC_URI))
	gosec -exclude-generated ./...

.PHONY: vuln
vuln:
	$(if $(GO_VULNCHECK), ,go install $(GO_VULNCHECK_URI))
	govulncheck ./...

.PHONY: verify
verify: lint sec vuln

.PHONY: mocks
mocks:
	$(if $(MOCK_GEN), ,go install $(MOCK_GEN_URI))
	@rm -rf $(MOCK_DESTINATION)
	mockgen -source=internal/ports/account_repository.go -destination=$(MOCK_DESTINATION)/account_repository_mock.go -package=mocks
	mockgen -source=internal/ports/transaction_repository.go -destination=$(MOCK_DESTINATION)/transaction_repository_mock.go -package=mocks
	mockgen -source=internal/file/filereader.go -destination=$(MOCK_DESTINATION)/filereader_mock.go -package=mocks
	mockgen -source=internal/notification/notifier.go -destination=$(MOCK_DESTINATION)/notifier_mock.go -package=mocks

.PHONY: test
test:
	go test -race -v -covermode=atomic -coverpkg=./... ./... -coverprofile=coverage.coverprofile
	go tool cover -html=coverage.coverprofile -o cover.html
	@echo "Coverage report generated at cover.html"

.PHONY: integration
integration:
	go test -tags=integration -v ./integration


