.PHONY: dep
dep:
	go mod tidy && go mod verify

.PHONY: mock
mock:
	@go generate $(PACKAGES)

.PHONY: test
test:
	go test -race ./...

.PHONY: cover
cover:
	go test -coverprofile=./profile.out ./...
	go tool cover -html=./profile.out

.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: build
build:
	go build -o bin/service ./cmd/

.PHONY: docker
docker:
	docker build -t service-sales:latest .

.PHONY: run
run:
	docker run --publish 8005:8005 service-sales:latest


