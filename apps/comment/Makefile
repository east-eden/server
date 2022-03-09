v ?= latest

GOVERSION=`go version`
BINARYVERSION=`git tag`
GITLASTLOG=$(shell git log --pretty=medium -1)
LDFLAGS=-ldflags '-w -X "main.BinaryVersion=${BINARYVERSION}" -X "main.GoVersion=${GOVERSION}" -X "main.GitLastLog=${GITLASTLOG}"'

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o comment main.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o comment.exe main.go

.PHONY: docker
docker: build
	docker build . -t comment:latest

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: push
push:
	# docker hub
	docker login -u hellodudu86 -p hello1986
	docker tag comment hellodudu86/comment:$(v)
	docker push hellodudu86/comment:$(v)

.PHONY: push_coding
push_coding:
	# coding hub
	docker login -u hellodudu86@gmail.com mmstudio-docker.pkg.coding.net
	docker tag comment mmstudio-docker.pkg.coding.net/blade/server/comment:$(v)
	docker push mmstudio-docker.pkg.coding.net/blade/server/comment:$(v)

	
.PHONY: push_github
push_github:
	# github packages
	docker tag comment ghcr.io/east-eden/server/comment:$(v)
	docker push ghcr.io/east-eden/server/comment:$(v)

.PHONY: clean
clean:
	make -C ../../ clean
