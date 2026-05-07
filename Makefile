.PHONY: verify test vet fmt-check smoke run

GOCACHE ?= /tmp/deutsch-tui-gocache

run:
	go run ./cmd/deutsch-tui

verify: fmt-check test vet smoke

fmt-check:
	@test -z "$$(gofmt -l cmd internal)"

test:
	GOCACHE=$(GOCACHE) go test ./...

vet:
	GOCACHE=$(GOCACHE) go vet ./...

smoke:
	./scripts/tui_smoke.sh

fmt:
	gofmt -w cmd internal

build:
	go build -o deutsch-tui-bin ./cmd/deutsch-tui

clean:
	rm -f deutsch-tui-bin
	rm -rf $(GOCACHE)

e2e:
	./scripts/verify.sh
