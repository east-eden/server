v ?= latest

GOVERSION=`go version`
BINARYVERSION=`git tag`
GITLASTLOG=$(shell git log --pretty=medium -1)
LDFLAGS=-ldflags '-w -X "main.BinaryVersion=${BINARYVERSION}" -X "main.GoVersion=${GOVERSION}" -X "main.GitLastLog=${GITLASTLOG}"'

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o mail main.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o mail.exe main.go

.PHONY: docker
docker: build
	docker build . -t mail:latest

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: push
push:
	@test -n "$$DOCKER_USER" || (echo "DOCKER_USER is required (Docker Hub account)"; exit 1)
	echo "$$DOCKER_PASS" | docker login -u "$$DOCKER_USER" --password-stdin
	docker tag mail $$DOCKER_USER/mail:$(v)
	docker push $$DOCKER_USER/mail:$(v)



	
.PHONY: push_github
push_github:
	# github packages
	docker tag mail ghcr.io/east-eden/server/mail:$(v)
	docker push ghcr.io/east-eden/server/mail:$(v)

.PHONY: clean
clean:
	make -C ../../ clean
