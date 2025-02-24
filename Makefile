ENTRYPOINT=cmd/server/main.go
GO_VERSION=1.22.12
APP_NAME=server

tidy:
	@go mod tidy -go=${GO_VERSION}

build: tidy
	@go build -o ./bin/${APP_NAME} ${ENTRYPOINT}

run: build
	@./bin/${APP_NAME} ${ARGS}

check:
	staticcheck ./...
	go vet ./...

test-all:
	go test -v ./...

lint:
	golangci-lint run --issues-exit-code 1 --print-issued-lines=true  ./...

docker-test:
	docker buildx build . \
		--build-arg APP_PATH=${ENTRYPOINT}