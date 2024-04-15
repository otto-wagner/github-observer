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
	go build -o app/github-observer .

build/run:
	./app/github-observer server

run:
	go run . server

docker/start:
	docker-compose up -d
