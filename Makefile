ci: setup test

setup:
	go mod tidy
	go mod verify

test:
	go test ./... -tags=unit,integration -coverprofile=coverage.out

test/unit:
	go test ./... -tags=unit -coverprofile=coverage.out

test/integration:
	go test ./... -tags=integration -coverprofile=coverage.out

mock: mock/setup mock/internal

mock/setup:
	go install github.com/vektra/mockery/v2@v2.43.2

mock/internal:
	cd ./internal && mockery --all

build:
	go build -o app/github-observer .

build/run:
	./app/github-observer server

run:
	go run . server

docker/start:
	docker-compose build && docker-compose up -d
