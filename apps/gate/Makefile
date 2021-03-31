v ?= latest

GOVERSION=`go version`
BINARYVERSION=`git tag`
GITLASTLOG=$(shell git log --pretty=medium -1)
LDFLAGS=-ldflags '-w -s -X "main.BinaryVersion=${BINARYVERSION}" -X "main.GoVersion=${GOVERSION}" -X "main.GitLastLog=${GITLASTLOG}"'

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o gate main.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o gate.exe main.go

.PHONY: docker
docker: build
	docker build . -t gate:latest

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: push
push:
	# docker hub
	docker login -u hellodudu86 -p hello1986
	docker tag gate hellodudu86/gate:$(v)
	docker push hellodudu86/gate:$(v)

.PHONY: push_coding
push_coding:
	# coding hub
	docker login -u hellodudu86@gmail.com mmstudio-docker.pkg.coding.net
	docker tag gate mmstudio-docker.pkg.coding.net/blade/server/gate:$(v)
	docker push mmstudio-docker.pkg.coding.net/blade/server/gate:$(v)

	
.PHONY: push_github
push_github:
	# github packages
	docker tag gate docker.pkg.github.com/east-eden/server/gate:$(v)
	docker push docker.pkg.github.com/east-eden/server/gate:$(v)

.PHONY: clean
clean:
	make -C ../../ clean
