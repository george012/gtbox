#!/bin/bash
ProductName=gtbox

echo "============================ ${ProductName} ============================"
echo "  1、发布 ${ProductName}"
echo "======================================================================"
read -p "$(echo -e "请输入版本号[例如；v0.0.1]")" versionStr

fileVersionLineNo=`cat ./version.go | grep -n "const VERSION =" | awk -F ":" '{print $1}'`

oldfileVersionStr=`cat ./version.go | grep -n "const VERSION =" | awk -F ":" '{print $2}'`

newVersionStr='const VERSION = ''"'$versionStr'"'
sed -i "" -e "${fileVersionLineNo}s/${oldfileVersionStr}/${newVersionStr}/g" ./version.go

REV_LIST=`git rev-list --tags --max-count=1`
APP_VERSION=`git describe --tags $REV_LIST`
APP_OLD_VERSION=${APP_VERSION%.*}.$((${APP_VERSION##*.}-2))

git add . && git commit -m "Update ${versionStr}"  && git tag $versionStr && git push && git push --tags && git tag -d $APP_OLD_VERSION