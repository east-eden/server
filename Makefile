v ?= latest

APPS = game gate mail rank comment combat client client_bots
APPS_win = $(addsuffix _win, $(APPS))
APPS_darwin = $(addsuffix _darwin, $(APPS))
MODS = code_generator
MODS_win = $(addsuffix _win, $(MODS))
MODS_darwin = $(addsuffix _darwin, $(MODS))
OUTPUT=build

GOVERSION=$(shell go version)
BINARYVERSION=$(shell git describe --tags)
GITLASTLOG=$(shell git log --pretty=format:'%h - %s (%cd) <%an>' -1)
GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-X "github.com/east-eden/server/version.BinaryVersion=${BINARYVERSION}" \
		-X "github.com/east-eden/server/version.GoVersion=${GOVERSION}" \
		-X "github.com/east-eden/server/version.GitLastLog=${GITLASTLOG}" \
		-w'

all: $(APPS) $(APPS_win) $(APPS_darwin)

$(APPS):
	GOOS=linux GOARCH=amd64 ${GOBUILD} -o $(OUTPUT)/$@ apps/$@/main.go
#	GOOS=linux GOARCH=amd64 GOAMD64=v3 ${GOBUILD} -o $(OUTPUT)/$@-v3 apps/$@/main.go

$(MODS):
	GOOS=linux GOARCH=amd64 ${GOBUILD} -o $(OUTPUT)/$@ cmd/$@/main.go

.PHONY: build $(APPS) $(MODS)
build: $(APPS) $(MODS)

$(APPS_win):
	GOOS=windows GOARCH=amd64 ${GOBUILD} -o $(OUTPUT)/$(subst _win,,$@).exe apps/$(subst _win,,$@)/main.go

$(MODS_win):
	GOOS=linux GOARCH=amd64 ${GOBUILD} -o $(OUTPUT)/$@ cmd/$@/main.go

.PHONY: build_win $(APPS_win) $(MODS_win)
build_win: $(APPS_win) $(MODS_win)

$(APPS_darwin):
	GOOS=darwin GOARCH=arm64 ${GOBUILD} -o $(OUTPUT)/$@ apps/$(subst _darwin,,$@)/main.go

$(MODS_darwin):
	GOOS=darwin GOARCH=arm64 ${GOBUILD} -o $(OUTPUT)/$@ cmd/$(subst _darwin,,$@)/main.go

.PHONY: build_win $(APPS_darwin) $(MODS_darwin)
build_darwin: $(APPS_darwin) $(MODS_darwin)

# .PHONY: proto
# proto:
# 	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/game/*.proto
# 	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/gate/*.proto
# 	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/combat/*.proto
# 	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/account/*.proto
# 	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/pubsub/*.proto

.PHONY: excel_gen
excel_gen:
	./build/code_generator

make_docker = sudo docker build --build-arg APPLICATION=$(1) -f Dockerfile.template -t $(1):latest .;
.PHONY: docker
docker: $(APPS)
	@ $(foreach app, $(APPS),$(call make_docker,$(app)))


.PHONY: ci_build_base
ci_build_base:
	docker build -f ci-building-base.Dockerfile -t ci-building-base .
	docker login -u hellodudu86@gmail.com mmstudio-docker.pkg.coding.net
	docker tag ci-building-base mmstudio-docker.pkg.coding.net/blade/server/ci-building-base
	docker push mmstudio-docker.pkg.coding.net/blade/server/ci-building-base

.PHONY: test
test:
	go test -v ./... -bench=. -benchmem -benchtime=100x
	# go test -v ./... -cover -coverprofile=test.out -bench=. -benchmem -benchtime=100x

.PHONY: test_html
test_html: test
	go tool cover -html=test.out

.PHONY: run
run:
	docker-compose up -d

.PHONY: push
push:
	make -C apps/game push
	make -C apps/gate push
	make -C apps/mail push
	make -C apps/rank push
	make -C apps/comment push
	# make -C apps/combat push
	# make -C apps/client_bots push

.PHONY: push_coding
push_coding:
	make -C apps/game push_coding
	make -C apps/gate push_coding
	make -C apps/mail push_coding
	make -C apps/rank push_coding
	make -C apps/comment push_coding
	# make -C apps/combat push_coding
	# make -C apps/client_bots push_coding

.PHONY: push_github
push_github:
	make -C apps/game push_github
	make -C apps/gate push_github
	make -C apps/mail push_github
	make -C apps/rank push_github
	make -C apps/comment push_github
	# make -C apps/combat push_github
	# make -C apps/client_bots push_github

.PHONY: clean
clean: stop
	rm -rf $(OUTPUT)/*
	docker rm -f $(shell docker ps -a -q)
	docker rmi -f $(shell docker images -a -q)

.PHONY: stop
stop:
	sudo docker-compose down
