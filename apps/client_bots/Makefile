v ?= latest

GOVERSION=`go version`
BINARYVERSION=`git tag`
GITLASTLOG=$(shell git log --pretty=medium -1)
LDFLAGS=-ldflags '-X "main.BinaryVersion=${BINARYVERSION}" -X "main.GoVersion=${GOVERSION}" -X "main.GitLastLog=${GITLASTLOG}"'

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o client_bots main.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o client_bots.exe main.go

.PHONY: docker
docker: build
	docker build . -t client_bots:latest

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: push
push:
	# docker hub
	docker login -u hellodudu86 -p hello1986
	docker tag client_bots hellodudu86/client_bots:$(v)
	docker push hellodudu86/client_bots:$(v)

.PHONY: push_coding
	# coding hub
	docker tag client_bots mmstudio-docker.pkg.coding.net/blade/client/client_bots:$(v)
	docker push mmstudio-docker.pkg.coding.net/blade/client/client_bots:$(v)

	
.PHONY: push_github
push_github:
	# github packages
	docker tag server_game docker.pkg.github.com/east-eden/server/client_bots:$(v)
	docker push docker.pkg.github.com/east-eden/server/client_bots:$(v)

.PHONY: run
run:
	go run main.go -config_file=config/client_bots/config.toml

.PHONY: proto
proto:
	make -C ../../ proto

.PHONY: clean
clean:
	make -C ../../ clean
