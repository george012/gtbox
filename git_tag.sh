#!/bin/bash

set -e

ProductName=gtbox
REPO_PFEX=george012/gtbox
VersionFile=./version.go

aVersion=`cat $VersionFile | grep -n "const VERSION =" | awk -F ":" '{print $2}'`
CurrentVersionString=`echo "${aVersion/'const VERSION = '/}" | sed 's/\"//g'`

versionStr=""

function to_run() {
    if [ -z "$1" ]; then
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
        return 0
    elif [ "$1" == "auto" ]; then
        baseStr=$(echo $CurrentVersionString | cut -d'.' -f1)     # Get the base version (v0)
        base=${baseStr//v/}                                       # Get the base version (0)
        major=$(echo $CurrentVersionString | cut -d'.' -f2)       # Get the major version (0)
        minor=$(echo $CurrentVersionString | cut -d'.' -f3)       # Get the minor version (1)

        minor=$((minor+1))                          # Increment the minor version
        if ((minor==100)); then                     # Check if minor version is 100
            minor=0                                 # Reset minor version to 0
            major=$((major+1))                      # Increment major version
        fi

        if ((major==100)); then                     # Check if major version is 100
            major=0                                 # Reset major version to 0
            base=$((base+1))                        # Increment base version
        fi

        versionStr="v${base}.${major}.${minor}"
        return 0
    else
      return 1
    fi
}

function git_handle() {
    fileVersionLineNo=`cat $VersionFile | grep -n "const VERSION =" | awk -F ":" '{print $1}'`

    oldfileVersionStr=`cat $VersionFile | grep -n "const VERSION =" | awk -F ":" '{print $2}'`

    newVersionStr='const VERSION = ''"'$versionStr'"'

    sed -i -e "${fileVersionLineNo}s/${oldfileVersionStr}/${newVersionStr}/g" $VersionFile

    ovs=${oldfileVersionStr#const VERSION = \"}
    APP_OLD_VERSION=${ovs%\"}
    PRE_DEL_VERSION=${APP_OLD_VERSION%.*}.$((${APP_OLD_VERSION##*.}-1))

    git add . \
    && git commit -m "Update ${versionStr}"  \
    && git tag $versionStr \
    && git push \
    && git push --tags \
    && git tag -f latest $versionStr \
    && git push -f origin latest \
    && git tag -d $PRE_DEL_VERSION
}


if to_run "$1"; then
    git_handle
else
    echo "Invalid argument"
fi


