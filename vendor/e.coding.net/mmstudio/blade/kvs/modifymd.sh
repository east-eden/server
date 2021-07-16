#!/bin/bash

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do SOURCE="$(readlink "$SOURCE")"; done
DIR_CURR="$(cd -P "$(dirname "$SOURCE")/" && pwd)"
release_version=$(cat "${DIR_CURR}"/package.json | jq -r '.version')
proc=`cat .versionrc| jq .productName|sed 's/\"//g'`

awk '{if($3 ~ /([0-9]{4}-[0-9]{2}-[0-9]{2})/ && date -d"$3")print "## ",$2,$3;else print $0}' CHANGELOG.md >new.md

i5ting_toc -f new.md
rm new.md
mv preview/new.html preview/release.html
sed -i ''  "s/i5ting_ztree_toc:new/${proc}/g" ./preview/release.html

sed -i ''  's/\"toc\//\"..\/chglog\/toc\//g' ./preview/release.html
sed -i ''  's/\"toc_/\"..\/chglog\/toc_/g' ./preview/release.html

expect -c "
spawn scp `pwd`/preview/release.html centurygame@10.0.84.230:/data/share/${proc}/release.html
expect {
\"*assword\" {set timeout 300; send \"centurygame\n\";}
\"yes/no\" {send \"yes\n\"; exp_continue;}
}
expect eof"
gob_passphrase=VLqZgVE8pBwJHtOG protokitgo release upload_file --remote_dir=software/${proc} --upload_file=`pwd`/preview/release.html
