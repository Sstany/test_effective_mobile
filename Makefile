.PHONY: lint
lint:
	golangci-lint --new-from-rev=main run

.PHONY: tests
tests:
	go test ./...

.PHONY: coverage
coverage:
	go test -cover ./internal/app/usecase

.PHONY: gen-coverage
gen-coverage:
	go test -cover ./internal/app/usecase -coverpkg=./internal/app/usecase -coverprofile ./coverage.out

.PHONY: show-coverage
show-coverage:
	go tool cover -func ./coverage.out

.PHONY: generate
generate:
	go generate ./...