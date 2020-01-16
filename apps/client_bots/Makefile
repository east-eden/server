v ?= latest

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o clientBotsBots main.go plugin.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o clientBots.exe main.go plugin.go

.PHONY: docker
docker: build
	docker build . -t yokai_server_clientBots:latest

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: push
push:
	docker tag yokai_server_clientBots hellodudu86/yokai_server_clientBots:$(v)
	docker push hellodudu86/yokai_server_clientBots:$(v)

.PHONY: run
run:
	go run main.go plugin.go -config_file=../../config/clientBots/config.toml

.PHONY: proto
proto:
	make -C ../../ proto

.PHONY: clean
clean:
	make -C ../../ clean
