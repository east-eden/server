.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o game main.go plugin.go

.PHONY: docker
docker: build
	docker build . -t yokai_server_game:latest

.PHONY: run
run:
	go run main.go plugin.go -config_file=../../config/config.toml


.PHONY: proto
proto:
	make -C ../../ proto
