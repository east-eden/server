v ?= latest

.PHONY: build
build: 
	env CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=linux GOARCH=amd64 go build -o crystal_tool crystal_tool.go

.PHONY: build_win
build_win: 
	env CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -o crystal_tool.exe crystal_tool.go

.PHONY: test
test:
	go test -v ./... -cover
