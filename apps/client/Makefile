v ?= latest

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o client main.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o client.exe main.go

.PHONY: docker
docker: build
	docker build . -t server_client:latest

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: push
push:
	docker tag server_client hellodudu86/server_client:$(v)
	docker push hellodudu86/server_client:$(v)

.PHONY: run
run:
	go run main.go -config_file=config/client/config.toml

.PHONY: proto
proto:
	make -C ../../ proto

.PHONY: clean
clean:
	make -C ../../ clean
