v ?= latest

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o battle main.go plugin.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o battle.exe main.go plugin.go

.PHONY: docker
docker: build
	docker build . -t yokai_server_battle:latest

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: push
push:
	docker tag yokai_server_battle hellodudu86/yokai_server_battle:$(v)
	docker push hellodudu86/yokai_server_battle:$(v)

.PHONY: run
run:
	go run main.go plugin.go -config_file=../../config/battle/config.toml

.PHONY: proto
proto:
	make -C ../../ proto

.PHONY: clean
clean:
	make -C ../../ clean
