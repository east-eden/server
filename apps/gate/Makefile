v ?= latest

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gate main.go plugin.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gate.exe main.go plugin.go

.PHONY: docker
docker: build
	docker build . -t yokai_server_gate:latest

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: push
push:
	docker tag yokai_server_gate hellodudu86/yokai_server_gate:$(v)
	docker tag yokai_server_gate mmstudio-docker.pkg.coding.net/blade/server/gate:$(v)
	docker push hellodudu86/yokai_server_gate:$(v)
	docker push mmstudio-docker.pkg.coding.net/blade/server/gate:$(v)

.PHONY: run
run:
	go run main.go plugin.go -config_file=../../config/gate/config.toml

.PHONY: proto
proto:
	make -C ../../ proto

.PHONY: clean
clean:
	make -C ../../ clean
