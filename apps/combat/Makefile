v ?= latest

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o combat main.go plugin.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o combat.exe main.go plugin.go

.PHONY: docker
docker: build
	docker build . -t yokai_server_combat:latest

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: push
push:
	# docker hub
	docker login -u hellodudu86 -p hello1986
	docker tag yokai_server_combat hellodudu86/yokai_server_combat:$(v)
	docker push hellodudu86/yokai_server_combat:$(v)

.PHONY: push_coding
	# coding hub
	docker tag yokai_server_combat mmstudio-docker.pkg.coding.net/blade/server/combat:$(v)
	docker push mmstudio-docker.pkg.coding.net/blade/server/combat:$(v)

	
.PHONY: push_github
push_github:
	# github packages
	cat ~/GH_TOKEN.txt | docker login docker.pkg.github.com -u hellodudu --password-stdin
	docker tag yokai_server_combat docker.pkg.github.com/yokaiio/yokai_server/yokai_server_combat:$(v)
	docker push yokai_server_combat docker.pkg.github.com/yokaiio/yokai_server/yokai_server_combat:$(v)

.PHONY: run
run:
	go run main.go plugin.go -config_file=config/combat/config.toml

.PHONY: proto
proto:
	make -C ../../ proto

.PHONY: clean
clean:
	make -C ../../ clean
