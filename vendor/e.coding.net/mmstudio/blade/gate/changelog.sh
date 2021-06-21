#!/bin/bash

SOURCE="${BASH_SOURCE[0]}"
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
cd "$DIR" || exit

if ! command -v conventional-changelog > /dev/null; then
  npm install -g conventional-changelog-cli
fi

conventional-changelog -p angular -i CHANGELOG.md -s
