v ?= latest

.PHONY: build
build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../excel/bin/code_generator main.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../../../excel/bin/code_generator.exe main.go

.PHONY: build_mac
build_mac: 
	env CGO_ENABLED=0 go build -o ../../../excel/bin/code_generator_mac main.go

.PHONY: build_all
build_all: build build_win build_mac


.PHONY: test
test:
	go test -v ./... -cover


.PHONY: proto
proto:
	make -C ../../ proto

.PHONY: clean
clean:
	make -C ../../ clean

