v ?= latest

.PHONY: build
build:
	make -C apps/game build
	make -C apps/gate build
	# make -C apps/combat build
	make -C apps/client build
	make -C apps/client_bots build
	make -C apps/code_generator build

.PHONY: build_win
build_win:
	make -C apps/game build_win
	make -C apps/gate build_win
	# make -C apps/combat build_win
	make -C apps/client build_win
	make -C apps/client_bots build_win

# .PHONY: proto
# proto:
# 	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/game/*.proto
# 	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/gate/*.proto
# 	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/combat/*.proto
# 	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/account/*.proto
# 	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/pubsub/*.proto

.PHONY: excel_gen
excel_gen:
	./apps/code_generator/code_generator

.PHONY: docker
docker:
	make -C apps/game docker
	make -C apps/gate docker
	# make -C apps/combat docker
	# make -C apps/client_bots docker

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
	make -C apps/combat push
	make -C apps/client_bots push

.PHONY: push_coding
push_coding:
	make -C apps/game push_coding
	make -C apps/gate push_coding
	# make -C apps/combat push_coding
	# make -C apps/client_bots push_coding

.PHONY: push_github
push_github:
	make -C apps/game push_github
	make -C apps/gate push_github
	make -C apps/combat push_github
	make -C apps/client_bots push_github

.PHONY: clean
clean:
	docker rm -f $(shell docker ps -a -q)
	docker rmi -f $(shell docker images -a -q)

.PHONY: stop
stop:
	docker-compose down
