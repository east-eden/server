v ?= latest

.PHONY: build
build:
	make -C apps/game build
	make -C apps/battle build
	make -C apps/client build

.PHONY: proto
proto:
	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/game/game.proto
	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/battle/battle.proto
	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/client/client.proto
	protoc -I=./proto --go_out=:${GOPATH}/src --micro_out=:${GOPATH}/src ./proto/pubsub/pubsub.proto

.PHONY: docker
docker:
	make -C apps/game docker
	make -C apps/battle docker

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: run
run:
	docker-compose up -d

.PHONY: push
push:
	make -C apps/game push
	make -C apps/battle push

.PHONY: clean
clean:
	mysql -uroot < config/sql/reset_game.sql < config/sql/reset_battle.sql
	docker rm -f $(shell docker ps -a -q)

.PHONY: stop
stop:
	docker-compose down
