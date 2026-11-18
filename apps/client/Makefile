v ?= latest

GOVERSION=`go version`
BINARYVERSION=`git tag`
GITLASTLOG=$(shell git log --pretty=medium -1)
LDFLAGS=-ldflags '-w -X "main.BinaryVersion=${BINARYVERSION}" -X "main.GoVersion=${GOVERSION}" -X "main.GitLastLog=${GITLASTLOG}"'

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o client main.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o client.exe main.go

.PHONY: build_mac
build_mac: 
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o client main.go

.PHONY: docker
docker: build
	docker build . -t server_client:latest

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: push
push:
	@test -n "$$DOCKER_USER" || (echo "DOCKER_USER is required (Docker Hub account)"; exit 1)
	echo "$$DOCKER_PASS" | docker login -u "$$DOCKER_USER" --password-stdin
	docker tag server_client $$DOCKER_USER/server_client:$(v)
	docker push $$DOCKER_USER/server_client:$(v)

.PHONY: run
run:
	go run main.go -config_file=config/client/config.toml

.PHONY: proto
proto:
	make -C ../../ proto

.PHONY: clean
clean:
	make -C ../../ clean
