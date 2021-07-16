#!/bin/bash

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do SOURCE="$(readlink "$SOURCE")"; done
DIR_CURR="$(cd -P "$(dirname "$SOURCE")/" && pwd)"
release_version=$(cat "${DIR_CURR}"/package.json | jq -r '.version')

run_test() {
   go test .
}

build() {
  if [  $? -ne 0 ];then
    echo "test error"
    exit 1
  fi

  rm -rf dist
  cd cmd/gate/

  BUILD_DATE=${BUILD_DATE:-$(date +%Y%m%d-%H:%M:%S)}
  revision=$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')
  branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'unknown')
  ldflags="-s -w -X e.coding.net/mmstudio/blade/golib/version.Version=$1"
  ldflags="$ldflags -X e.coding.net/mmstudio/blade/golib/version.Revision=$revision"
  ldflags="$ldflags -X e.coding.net/mmstudio/blade/golib/version.Branch=$branch"
  ldflags="$ldflags -X e.coding.net/mmstudio/blade/golib/version.BuildDate=$BUILD_DATE"

  CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -gcflags=all=-trimpath=$DIR_CURR -asmflags=all=-trimpath=$DIR_CURR  -ldflags "${ldflags}" -o ../../dist/osx/gate main.go
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags=all=-trimpath=$DIR_CURR -asmflags=all=-trimpath=$DIR_CURR  -ldflags "${ldflags}" -o ../../dist/linux/gate main.go
}

case "$1" in
  test)
    run_test
  ;;
  build)
	  build $release_version
	;;
esac
