#!/bin/sh
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do SOURCE="$(readlink "$SOURCE")"; done
DIR_CURR="$(cd -P "$(dirname "$SOURCE")/" && pwd)"
release_version=$(cat "${DIR_CURR}"/package.json | jq -r '.version')
proc=`cat .versionrc| jq .productName|sed 's/\"//g'`
if [ \"$1\" == \"stable\" ]
 then
	result=`echo $release_version |grep beta`
	if [[ "$result" != "" ]]
	then
		newversion=`echo $release_version |awk -F'-' '{print $1}'`
		line=`grep -n -v "${newversion}-beta" CHANGELOG.md|grep -E "([0-9]{4}-[0-9]{2}-[0-9]{2})"|head -1|awk -F':' '{print $1}'`;
		preversion=$(($line - 1))
		sed -i '' "4,${preversion}d" CHANGELOG.md
		git tag |grep "${newversion}-beta" |xargs git tag -d
	else
		newversion=`echo $release_version |awk -F'.' '{v=$1;for(i=2;i<NF;i++)v=v"."$i;print v"."$NF+1}'`
	fi
else
	result=`echo $release_version |grep beta`
	if [[ "$result" != "" ]]
	then
		newversion=`echo $release_version |awk -F'.' '{v=$1;for(i=2;i<NF;i++)v=v"."$i;print v"."$NF+1}'`
	else
		newversion=`echo $release_version |awk -F'.' '{v=$1;for(i=2;i<NF;i++)v=v"."$i;print v"."$NF+1"-beta.1"}'`
	fi
fi
echo $newversion;
npm run release -- --release-as $newversion;
echo "Release URL:https://zhongtai.s3.amazonaws.com/software/${proc}/release.html"
