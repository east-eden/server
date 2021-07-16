DIR_CURR="$(cd -P "$(dirname "$SOURCE")/" && pwd)"

rm -rf preview/*
echo "### Download URL" >zip.txt
echo '' >>zip.txt
echo " * [Mac](https://zhongtai.s3.amazonaws.com/software/gate/protokitgo_osx_$release_version.zip)  /" >>zip.txt
#echo "[Windows](https://zhongtai.s3.amazonaws.com/software/gate/protokitgo_win_$release_version.zip)  /">>zip.txt
echo "[Linux](https://zhongtai.s3.amazonaws.com/software/gate/protokitgo_linux_$release_version.zip)">>zip.txt
sed -i '' '5 r zip.txt' CHANGELOG.md
rm -f zip.txt

cat CHANGELOG.md >tmp.md

awk '{if($2~/funplus\/gate\/compare/)print "## ",$2,$3;else print $0}' tmp.md >new.md
rm tmp.md

i5ting_toc -f new.md
rm new.md
sed -i ''  's/is_auto_number:true,/is_auto_number:false,/g' ./preview/toc_conf.js
mv preview/new.html preview/release.html
sed -i ''  's/i5ting_ztree_toc:new/gate/g' ./preview/release.html


gob_passphrase=VLqZgVE8pBwJHtOG protokitgo release upload_file --remote_dir=software/gate --upload_file=$DIR_CURR/preview/release.html