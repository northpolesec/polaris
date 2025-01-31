default: fmt build

build:
	go build -v ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

run:
	docker compose up --build

.PHONY: fmt test build run
