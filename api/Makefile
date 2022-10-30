.PHONY: install test build serve clean pack deploy ship

TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)

export TAG

install:
	go get .

test: install
	go test ./...

build: install
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(TAG)" -o ./bin/mail-api .

serve: build
	./bin/mail-api serve grpc

clean:
	rm -f ./bin/mail-api

dev:
	make build
	docker build -t localhost:5000/mail-api .
	docker push localhost:5000/mail-api
