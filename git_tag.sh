#!/bin/bash

set -e

ProductName=$(grep ProjectName ./config/config.go | awk -F '"' '{print $2}' | sed 's/\"//g')
Product_version_key="ProjectVersion"
REPO_PFEX=george012/$ProductName
VersionFile=./config/config.go

CURRENT_VERSION=$(grep ${Product_version_key} $VersionFile | awk -F '"' '{print $2}' | sed 's/\"//g')

NEXT_VERSION=""

OS_TYPE="Unknown"
GetOSType() {
    uNames=`uname -s`
    osName=${uNames: 0: 4}
    if [ "$osName" == "Darw" ] # Darwin
    then
        OS_TYPE="Darwin"
    elif [ "$osName" == "Linu" ] # Linux
    then
        OS_TYPE="Linux"
    elif [ "$osName" == "MING" ] # MINGW, windows, git-bash
    then
        OS_TYPE="Windows"
    else
        OS_TYPE="Unknown"
    fi
}
GetOSType

function to_run() {
    if [ -z "$1" ]; then
        baseStr=$(echo ${CURRENT_VERSION} | cut -d'.' -f1)     # Get the base version (v0)
        base=${baseStr//v/}                                       # Get the base version (0)
        major=$(echo ${CURRENT_VERSION} | cut -d'.' -f2)       # Get the major version (0)
        minor=$(echo ${CURRENT_VERSION}| cut -d'.' -f3)       # Get the minor version (1)

        minor=$((minor+1))                          # Increment the minor version
        if ((minor==1000)); then                     # Check if minor version is 100
            minor=0                                 # Reset minor version to 0
            major=$((major+1))                      # Increment major version
        fi

        if ((major==1000)); then                     # Check if major version is 100
            major=0                                 # Reset major version to 0
            base=$((base+1))                        # Increment base version
        fi

        NEXT_VERSION="v${base}.${major}.${minor}"
        return 0
    elif [ "$1" == "custom" ]; then
        echo "============================ ${ProductName} ============================"
        echo "  1、发布 [-${ProductName}-]"
        echo "  当前版本[-${CURRENT_VERSION}-]"
        echo "======================================================================"
        read -p "$(echo -e "请输入版本号[例如；v0.0.1]")" inputString
        if [[ "$inputString" =~ ^v.* ]]; then
            NEXT_VERSION=${inputString}
        else
            NEXT_VERSION=v${inputString}
        fi
        return 0
    else
        return 1
    fi
}

function get_pre_del_version_no {
    local v_str=$1
    baseStr=$(echo $v_str | cut -d'.' -f1)     # Get the base version (v0)
    base=${baseStr//v/}                                       # Get the base version (0)
    major=$(echo $v_str | cut -d'.' -f2)       # Get the major version (0)
    minor=$(echo $v_str | cut -d'.' -f3)       # Get the minor version (1)

    if ((minor>0)); then                      # Check if minor version is more than 0
        minor=$((minor-1))                     # Decrement the minor version
    else
        minor=999                              # Reset minor version to 99
        if ((major>0)); then                   # Check if major version is more than 0
            major=$((major-1))                 # Decrement major version
        else
            major=999                           # Reset major version to 99
            if ((base>0)); then                # Check if base version is more than 0
                base=$((base-1))               # Decrement base version
            else
                echo "Error: Version cannot be decremented."
                exit 1
            fi
        fi
    fi

    pre_v_no="${base}.${major}.${minor}"
    echo $pre_v_no
}

function git_handle_ready() {
    echo "Current Version With "${CURRENT_VERSION}
    echo "Next Version With "${NEXT_VERSION}

    sed -i -e "s/\(${Product_version_key}[[:space:]]*=[[:space:]]*\"\)${CURRENT_VERSION}\"/\1${NEXT_VERSION}\"/" $VersionFile

    if [[ $OS_TYPE == "Darwin" ]]; then
        echo "rm darwin cache"
        rm -rf $VersionFile"-e"
    fi
}

function git_handle_push() {
    local current_version_no=${CURRENT_VERSION//v/}
    local next_version_no=${NEXT_VERSION//v/}
    local pre_del_version_no=$(get_pre_del_version_no "$current_version_no")
    echo "Pre Del Version With v"${pre_del_version_no}

    git add . \
    && git commit -m "Release v${next_version_no}_$(date -u +"%Y-%m-%d_%H:%M:%S")"_"UTC" \
    && git tag v${next_version_no} \
    && git tag -f latest v${next_version_no}

    for remote in $(git remote)
    do
        echo "Pushing to ${remote}..."
        git push --delete ${remote} latest \
        && git push ${remote} \
        && git push ${remote} v${next_version_no} \
        && git push ${remote} latest
    done
    git tag -d v${pre_del_version_no}
}

handle_input(){
    if [[ $1 == "-get_pre_del_tag_name" ]]; then
        pre_tag=$(get_pre_del_version_no "${CURRENT_VERSION}")
        echo "Pre Del Tag With " "$pre_tag"
    elif [ -z "$1" ] || [ "$1" == "auto" ]; then

        if to_run "$1"; then
            git_handle_ready
            git_handle_push
            echo "Complated"
        else
            echo "Invalid argument normal"
        fi
    else
        echo "Invalid argument"
    fi
}

handle_input "$@"
