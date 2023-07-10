#!/bin/bash
set -e

PAT=$1

ProductName=gtbox
REPO_PFEX=george012/gtbox
VersionFile=./version.go

aVersion=`cat $VersionFile | grep -n "const VERSION =" | awk -F ":" '{print $2}'`
CurrentVersionString=`echo "${aVersion/'const VERSION = '/}" | sed 's/\"//g'`

versionStr=""

if [ -z "$2" ]; then
    echo "============================ ${ProductName} ============================"
    echo "  1、发布 [-${ProductName}-]"
    echo "  当前版本[-${CurrentVersionString}-]"
    echo "======================================================================"
    read -p "$(echo -e "请输入版本号[例如；v0.0.1]")" inputString
    if [[ "$inputString" =~ ^v.* ]]; then
        versionStr=${inputString}
    else
        versionStr=v${inputString}
    fi
elif [ "$2" == "auto" ]; then
    base=$(echo $CurrentVersionString | cut -d'.' -f1)      # Get the base version (v0)
    major=$(echo $CurrentVersionString | cut -d'.' -f2)      # Get the major version (0)
    minor=$(echo $CurrentVersionString | cut -d'.' -f3)      # Get the minor version (1)

    minor=$((minor+1))                          # Increment the minor version
    if ((minor==100)); then                     # Check if minor version is 100
        minor=0                                 # Reset minor version to 0
        major=$((major+1))                      # Increment major version
    fi

    if ((major==100)); then                     # Check if major version is 100
        major=0                                 # Reset major version to 0
        base="${base}1"                         # Increment base version
    fi

    versionStr="${base}.${major}.${minor}"
fi



fileVersionLineNo=`cat $VersionFile | grep -n "const VERSION =" | awk -F ":" '{print $1}'`

oldfileVersionStr=`cat $VersionFile | grep -n "const VERSION =" | awk -F ":" '{print $2}'`

newVersionStr='const VERSION = ''"'$versionStr'"'
sed -i "" -e "${fileVersionLineNo}s/${oldfileVersionStr}/${newVersionStr}/g" $VersionFile

ovs=${oldfileVersionStr#const VERSION = \"}
APP_OLD_VERSION=${ovs%\"}
PRE_DEL_VERSION=${APP_OLD_VERSION%.*}.$((${APP_OLD_VERSION##*.}-1))

if [[ -z "$PAT" ]]; then
    git add . \
    && git commit -m "Update ${versionStr}"  \
    && git tag $versionStr \
    && git push \
    && git push --tags \
    && git tag -f latest $versionStr \
    && git push -f origin latest \
    && git tag -d $PRE_DEL_VERSION
else
    git add . \
    && git commit -m "Update ${versionStr}"  \
    && git tag $versionStr \
    && git push https://$PAT@github.com/${REPO_PFEX}.git \
    && git push --tags https://$PAT@github.com/${REPO_PFEX}.git \
    && git tag -f latest $versionStr \
    && git push -f https://$PAT@github.com/${REPO_PFEX}.git origin latest \
    && git tag -d $PRE_DEL_VERSION
fi