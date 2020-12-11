v ?= latest

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chat main.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o chat.exe main.go

.PHONY: docker
docker: build
	docker build . -t server_chat:latest

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: push
push:
	docker tag server_chat hellodudu86/server_chat:$(v)
	docker tag server_chat mmstudio-docker.pkg.coding.net/blade/server/chat:$(v)
	docker push hellodudu86/server_chat:$(v)
	docker push mmstudio-docker.pkg.coding.net/blade/server/chat:$(v)

.PHONY: run
run:
	go run main.go -config_file=config/chat/config.toml

.PHONY: proto
proto:
	make -C ../../ proto

.PHONY: clean
clean:
	make -C ../../ clean
