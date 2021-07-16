TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: cover

tools:
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	@go tool goimports 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/goimports; \
	fi
	@go tool optiongen 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u github.com/timestee/optiongen/cmd/optiongen; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

test:
	go list  $(TEST) | xargs -t -n4 go test $(TESTARGS)  -timeout=2m -parallel=4

race:
	go test  -race $(TEST) $(TESTARGS)


cover:tools
	go test $(TEST) -coverprofile=coverage.out.tmp
	cat coverage.out.tmp | grep -v "gen.go"|grep -v "gen_"| grep -v vendor > coverage.out
	go tool cover -html=coverage.out
	rm coverage.out
	rm coverage.out.tmp


generate: tools
	go generate ./...


snore:
	-golangci-lint run --out-format checkstyle > ./golangci-lint.out
	go test $(TEST) -coverprofile=coverage.out.tmp
	cat coverage.out.tmp | grep -v "gen.go"|grep -v "gen_"|grep -v "/pkg" | grep -v vendor | grep -v test_perf | grep -v samples  | grep -v internal > coverage.out
	sonar-scanner
	rm coverage.out
	rm coverage.out.tmp
	rm golangci-lint.out

.NOTPARALLEL:
.PHONY: cover fmt generate test tools race cover snore

stable:
	sh  $(CURDIR)/release.sh stable

beta:
	sh  $(CURDIR)/release.sh beta


stable:
	sh  $(CURDIR)/release.sh stable

beta:
	sh  $(CURDIR)/release.sh beta

