.PHONY: install test build serve clean pack deploy ship

TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)

export TAG

install:
	go get .

test: install
	go test ./...

build: install
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(TAG)" -o ./bin/send-mail-consumer .

serve: build
	./bin/send-mail-consumer serve grpc

clean:
	rm -f ./bin/send-mail-consumer

dev:
	make build
	docker build -t localhost:5000/send-mail-consumer .
	docker push localhost:5000/send-mail-consumer
