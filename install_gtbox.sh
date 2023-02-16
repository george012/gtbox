#!/bin/bash
ProductName=gtbox
#CustomLibs="gtgo,otherlib1,otherlib2"
CustomLibs=$(ls -l ./libs|awk '/^d/ {print $NF}')
echo $CustomLibs

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

removeCache() {
    rm -rf ./version.go
    rm -rf ./install_${ProductName}.sh
}

install() {
    GetOSType
    echo ${OSTYPE}
    if [ ${OSTYPE} == "Windows" ]
    then
        ago_path_dir=`echo "${GOPATH/':\\'/'/'}" | sed 's/\"//g'`
        complate_gopath_dir='/'`echo "${ago_path_dir}" | tr A-Z a-z`
        find $complate_gopath_dir/pkg/mod/github.com/george012  -name "${ProductName}@*" -exec rm -rf {} \;
    else
        find ${GOPATH}/pkg/mod/github.com/george012  -name "${ProductName}@*" -exec rm -rf {} \;
    fi

    go get -u github.com/george012/${ProductName}@latest

    rm -rf ./version.go
    wget --no-check-certificate https://raw.githubusercontent.com/george012/${ProductName}/master/version.go
    aVersion=`cat ./version.go | grep -n "const VERSION =" | awk -F ":" '{print $2}'`
    aVersionNo=`echo "${aVersion/'const VERSION = '/}" | sed 's/\"//g'`

    for libName in ${CustomLibs}
    do
        if [ ${OSTYPE} == "Darwin" ] # Darwin
        then
            srcPWD=`pwd`
            sudo rm -rf /usr/local/lib/lib${libName}_arm64.dylib
            sudo rm -rf /usr/local/lib/lib${libName}.dylib
    #        cd ${GOPATH}/pkg/mod/github.com/george012/gtbox@v${aVersionNo} && /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/install_name_tool -add_rpath ../gtbox@v${aVersionNo} ${produckName} && cd ${srcPWD}
            sudo ln -s ${GOPATH}/pkg/mod/github.com/george012/${ProductName}@v${aVersionNo}/libs/${libName}/lib${libName}.dylib /usr/local/lib/lib${libName}.dylib
            sudo ln -s /usr/local/lib/lib${libName}.dylib /usr/local/lib/lib${libName}_arm64.dylib

        elif [ ${OSTYPE} == "Linux" ] # Linux
        then
            rm -rf /lib64/lib${libName}.so
            ln -s ${GOPATH}/pkg/mod/github.com/george012/${ProductName}@v${aVersionNo}/libs/${libName}/lib${libName}.so /lib64/lib${libName}.so && ldconfig
        elif [ ${OSTYPE} == "Windows" ] # MINGW, windows, git-bash
        then
            rm -rf /c/Windows/System32/lib${libName}.dll
            ln -s ${GOPATH}/pkg/mod/github.com/george012/${ProductName}@v${aVersionNo}/libs/${libName}/${libName}.dll /c/Windows/System32/${libName}.dll
        else
            echo ${OSTYPE}
        fi
    done



    removeCache
}

uninstall() {
    GetOSType
    echo ${OSTYPE}
    removeCache
    for libName in ${CustomLibs}
    do
        if [ ${OSTYPE} == "Darwin" ] # Darwin
        then
            rm -rf /usr/local/lib/lib"${libName}"_arm64.dylib
            rm -rf /usr/local/lib/lib"${libName}".dylib
        elif [ ${OSTYPE} == "Linux" ] # Linux
        then
            rm -rf /lib64/lib"${libName}".so
        elif [ ${OSTYPE} == "Windows" ] # MINGW, windows, git-bash
        then
            rm -rf c:/Windows/System32/"${libName}".dll
        else
            echo ${OSTYPE}
        fi
    done


    removeCache

    find ${GOPATH}/pkg/mod/github.com/george012  -name "${ProductName}@*" -exec rm -rf {} \;
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
