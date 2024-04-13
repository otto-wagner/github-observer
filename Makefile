ci: setup test

setup:
	go mod tidy
	go mod verify

test: test/unit test/integration

test/unit:
	go test ./... -tags=unit

test/integration:
	go test ./... -tags=integration

mock: mock/setup mock/internal mock/pkg

mock/setup:
	go install github.com/vektra/mockery/v2@v2.40.3

mock/internal:
	cd ./internal && mockery --all

mock/pkg:
	cd ./pkg && mockery --all

build:
	go build -o bin/github-listener .

build/run:
	./bin/github-listener server

run:
	go run . server
