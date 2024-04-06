ci: setup test

setup:
	go mod tidy
	go mod verify

test: test/unit test/integration

test/unit:
	go test ./... -tags=unit

test/integration:
	go test ./... -tags=integration
