GO_BIN ?= go

all: lint

update:
	$(GO_BIN) get -u
	$(GO_BIN) get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GO_BIN) get -u github.com/mgechev/revive
	$(GO_BIN) mod tidy

lint:
	golangci-lint run ./...
	revive -config revive.toml -exclude headers.go ./...
