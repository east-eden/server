v ?= latest

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gate main.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o gate.exe main.go

.PHONY: docker
docker: build
	docker build . -t server_gate:latest

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: push
push:
	# docker hub
	docker login -u hellodudu86 -p hello1986
	docker tag server_gate hellodudu86/server_gate:$(v)
	docker push hellodudu86/server_gate:$(v)

.PHONY: push_coding
	# coding hub
	docker tag server_gate mmstudio-docker.pkg.coding.net/blade/server/gate:$(v)
	docker push mmstudio-docker.pkg.coding.net/blade/server/gate:$(v)

	
.PHONY: push_github
push_github:
	# github packages
	docker tag server_gate docker.pkg.github.com/east-eden/server/server_gate:$(v)
	docker push docker.pkg.github.com/east-eden/server/server_gate:$(v)


.PHONY: run
run:
	go run main.go -config_file=config/gate/config.toml

.PHONY: proto
proto:
	make -C ../../ proto

.PHONY: clean
clean:
	make -C ../../ clean
