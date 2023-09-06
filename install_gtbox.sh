#!/bin/bash

set -e

ProductName=gtbox
REPO_PFEX=george012/${ProductName}

OSTYPE="Unknown"
GetOSType(){
    uNames=`uname -s`
    osName=${uNames: 0: 4}
    if [ "$osName" == "Darw" ] # Darwin
    then
        OSTYPE="Darwin"
    elif [ "$osName" == "Linu" ] # Linux
    then
        OSTYPE="Linux"
    elif [ "$osName" == "MING" ] # MINGW, windows, git-bash
    then
        OSTYPE="Windows"
    else
        OSTYPE="Unknown"
    fi
}
GetOSType

removeCache() {
    rm -rf ./${ProductName}_config.go
    rm -rf ./install_${ProductName}.sh
}

parse_json(){
    echo "${1//\"/}" | tr -d '\n' | tr -d '\r' | sed "s/.*$2:\([^,}]*\).*/\1/"
}

get_repo_latest_version(){
    local REMOTE_REPO_VERSION=""
    local LATEST_RELEASE_INFO=$(curl --silent https://api.github.com/repos/${REPO_PFEX}/releases/latest)
    if ! echo "$LATEST_RELEASE_INFO" | grep -q "Not Found"; then
        REMOTE_REPO_VERSION=$(parse_json "$LATEST_RELEASE_INFO" "tag_name")
    else
      return 1
    fi
    echo $REMOTE_REPO_VERSION | tr -d '\r\n'
    return 0
}

create_symlink() {
    local alibName=$1
    local aVersionStr=$2
    local prefix="lib"
    local libPath=${complate_gopath_dir}/pkg/mod/github.com/george012/${ProductName}@${aVersionStr}/libs/${alibName}

    case ${OSTYPE} in
        "Darwin"|"Linux")
            # 如果 alibName 不是以 "lib" 开头，则添加 "lib" 前缀
            [[ ${alibName} == lib* ]] || alibName="${prefix}${alibName}"

            if [ "${OSTYPE}" == "Darwin" ]; then
                sudo ln -s ${libPath}/${alibName}.dylib /usr/local/lib/${alibName}.dylib
                sudo ln -s /usr/local/lib/${alibName}.dylib /usr/local/lib/${alibName}_arm64.dylib
            else
                ln -s ${libPath}/${alibName}.so /lib64/${alibName}.so && ldconfig
            fi
            ;;
        "Windows")
            [[ ${alibName} != lib* ]] || alibName="${alibName#lib}"
            ln -s ${libPath}/${alibName}.dll /c/Windows/System32/${alibName}.dll
            ;;
        *)
            echo ${OSTYPE}
            ;;
    esac
}


install() {
    echo ${OSTYPE}

    complate_gopath_dir=${GOPATH}
    if [ ${OSTYPE} == "Windows" ]
    then
        ago_path_dir=`echo "${GOPATH/':\\'/'/'}" | sed 's/\"//g'`
        complate_gopath_dir='/'`echo "${ago_path_dir}" | tr A-Z a-z`
    fi

    find ${complate_gopath_dir}/pkg/mod/github.com/george012 -depth -name "${ProductName}@*" -exec sudo rm -rf {} \;

    last_repo_version=$(get_repo_latest_version)

    go get -u github.com/george012/${ProductName}@${last_repo_version} \
    && {
        CustomLibs=$(ls -l ${complate_gopath_dir}/pkg/mod/github.com/george012/gtbox@${last_repo_version}/libs |awk '/^d/ {print $NF}') \
        && for alibName in ${CustomLibs}
        do
            create_symlink ${alibName} ${last_repo_version}
        done
    }

    removeCache
}

uninstall() {
    complate_gopath_dir=${GOPATH}

    # 找到所有版本的库并删除
    find ${complate_gopath_dir}/pkg/mod/github.com/george012/${ProductName}@* -type d -exec rm -rf {} \;

    # 删除所有自定义库
    CustomLibs=$(ls -l ${complate_gopath_dir}/pkg/mod/github.com/george012/${ProductName}/libs |awk '/^d/ {print $NF}')

    for libName in ${CustomLibs}
    do
        if [ ${OSTYPE} == "Darwin" ] # Darwin
        then
            rm -rf /usr/local/lib/lib${libName}_arm64.dylib
            rm -rf /usr/local/lib/lib${libName}.dylib
        elif [ ${OSTYPE} == "Linux" ] # Linux
        then
            rm -rf /lib64/lib${libName}.so
        elif [ ${OSTYPE} == "Windows" ] # MINGW, windows, git-bash
        then
            ago_path_dir=`echo "${GOPATH/':\\'/'/'}" | sed 's/\"//g'`
            complate_gopath_dir='/'`echo "${ago_path_dir}" | tr A-Z a-z`
            rm -rf /c/Windows/System32/${libName}.dll
        else
            echo ${OSTYPE}
        fi
    done

    removeCache
}

echo "============================ ${ProductName} ============================"
echo "  1、安装 ${ProductName}"
echo "  2、卸载 ${ProductName}"
echo "======================================================================"
read -p "$(echo -e "请选择[1-2]：")" choose
case $choose in
1)
    install
    ;;
2)
    uninstall
    ;;
*)
    echo "输入错误，请重新输入！"
    ;;
esac
