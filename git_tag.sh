#!/bin/bash
ProductName=gtbox

echo "============================ ${ProductName} ============================"
echo "  1、安装 ${ProductName}"
echo "======================================================================"
read -p "$(echo -e "请输入版本号[例如；v0.0.1]")" newVersionStr

fileVersionLineNo=`cat ./version.go | grep -n "const VERSION =" | awk -F ":" '{print $1}'`

oldfileVersionStr=`cat ./version.go | grep -n "const VERSION =" | awk -F ":" '{print $2}'`

sed -i "" -e "${fileVersionLineNo}s/${oldfileVersionStr}/${newVersionStr}/g" ./version.go

git add . && git commit -m "Update"  && git tag $newVersionStr && git push && git push --tags
