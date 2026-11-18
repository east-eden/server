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

# proto 源文件维护在独立仓库 github.com/east-eden/proto，
# 生成结果落在本仓库的 proto/ 目录，因此这里只做调用转发。
PROTO_DIR ?= ../proto

.PHONY: proto
proto:
	@test -d $(PROTO_DIR) || (echo "proto repo not found: $(PROTO_DIR)"; echo "set PROTO_DIR=/path/to/east-eden/proto"; exit 1)
	$(MAKE) -C $(PROTO_DIR) proto

.PHONY: excel_gen
excel_gen:
	./build/code_generator

# 生成本地开发用自签名证书（私钥不入库）
.PHONY: cert
cert:
	@mkdir -p config/cert
	openssl req -x509 -newkey rsa:2048 -nodes -days 3650 \
		-keyout config/cert/server.key -out config/cert/server.crt \
		-subj "/CN=localhost"
	@echo "generated config/cert/server.crt and config/cert/server.key"

make_docker = sudo docker build --build-arg APPLICATION=$(1) -f Dockerfile.template -t $(1):latest .;
.PHONY: docker
docker: $(APPS)
	@ $(foreach app, $(APPS),$(call make_docker,$(app)))


.PHONY: ci_build_base
ci_build_base:
	docker build -f ci-building-base.Dockerfile -t ci-building-base .

.PHONY: test
# 单元测试：不依赖外部中间件
test:
	go test ./...

.PHONY: test_verbose
test_verbose:
	go test -v ./...

.PHONY: test_bench
test_bench:
	go test -v ./... -bench=. -benchmem -benchtime=100x

.PHONY: test_integration
# 集成测试：需要本地 MongoDB（见 docker-compose.yaml）
test_integration:
	go test -tags=integration -v ./...

.PHONY: lint
lint:
	go vet ./...

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
