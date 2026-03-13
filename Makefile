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

generate/mocks:
	mockery

generate/certificate:
	./scripts/generate-certificate.sh

build:
	go build -o ./app/github-observer ./cmd/observer/main.go

build/run:
	./app/github-observer server run

run:
	go run ./cmd/observer/main.go server

docker/start:
	docker-compose build && docker-compose up -d

release: build-release
	docker push ottowagner/observer:$(VERSION) &
	docker push ottowagner/observer:latest

build-release:
	docker build -f deployments/observer/Dockerfile -t ottowagner/observer:$(VERSION) -t ottowagner/observer:latest .

