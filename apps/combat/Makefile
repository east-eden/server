v ?= latest

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o combat main.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o combat.exe main.go

.PHONY: docker
docker: build
	docker build . -t server_combat:latest

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: push
push:
	# docker hub
	docker login -u hellodudu86 -p hello1986
	docker tag server_combat hellodudu86/server_combat:$(v)
	docker push hellodudu86/server_combat:$(v)

.PHONY: push_coding
	# coding hub
	docker tag server_combat mmstudio-docker.pkg.coding.net/blade/server/combat:$(v)
	docker push mmstudio-docker.pkg.coding.net/blade/server/combat:$(v)

	
.PHONY: push_github
push_github:
	# github packages
	docker tag server_combat docker.pkg.github.com/east-eden/server/server_combat:$(v)
	docker push docker.pkg.github.com/east-eden/server/server_combat:$(v)

.PHONY: run
run:
	go run main.go -config_file=config/combat/config.toml

.PHONY: proto
proto:
	make -C ../../ proto

.PHONY: clean
clean:
	make -C ../../ clean
