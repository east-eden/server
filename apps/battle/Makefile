.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o battle main.go plugin.go

.PHONY: docker
docker: build
	docker build . -t yokai_server_battle:latest

.PHONY: run
run:
	go run main.go plugin.go -config_file=../../config/battle/config.toml

.PHONY: proto
proto:
	make -C ../../ proto

.PHONY: clean
clean:
	make -C ../../ clean
