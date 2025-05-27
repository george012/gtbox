<!-- TOC -->

- [1. 使用说明](#1-使用说明)
    - [1.1. 关于测试](#11-关于测试)
- [2. 使用](#2-使用)
    - [2.1. 静态库版本](#21-静态库版本)
    - [2.2. 动态库版本---停止维护](#22-动态库版本---停止维护)
    - [2.3. mac 编译 linux](#23-mac-编译-linux)
- [3. 更新记录](#3-更新记录)
- [4. 必要支持](#4-必要支持)
- [5. 功能](#5-功能)
- [6. 感谢](#6-感谢)

<!-- /TOC -->

# 1. 使用说明
* <font color=red>只保留2个可运行版本，方便维护</font>
* <font color=red>尝试改动态库为静态库 </font>
* 必须安装了`git`
* 项目必须用`go mod`自动管理依赖
* 必须：`go version` `>=` `1.18`

## 1.1. 关于测试 
* `go test -v -run ./...`

# 2. 使用
## 2.1. 静态库版本
```
go get -u github.com/george012/gtbox@latest
```
## 2.2. 动态库版本---停止维护
* <font color=red>在任意`golang`项目根目录下使用`terminal`执行如下命令</font>
```
wget --no-check-certificate https://raw.githubusercontent.com/george012/gtbox/master/install_gtbox.sh && chmod a+x ./install_gtbox.sh && ./install_gtbox.sh
```

## 2.3. mac 编译 linux
```
brew install filosottile/musl-cross/musl-cross

设置如下环境变量 到 ~/.bash_profile 或者 ~/.zshrc
# Musl-cross 环境
export MUSL_CROSS_ROOT=$(brew --prefix musl-cross)
export PATH=$MUSL_CROSS_ROOT/bin:$PATH


```

# 3. 更新记录
* [ChangeList](./ChangeList.md)

# 4. 必要支持
*   CGO支持
*   MAC安装最新版本Xcode及Command Line Tools


# 5. 功能
- ✅ CGO支持
- ✅ 自定义加、解密
- ✅ 简单的 SSH Client
- ✅ 简易 SCP 工具
- ✅ 简单的 HTTP Client
- ✅ 简单的 ORM 封装
- ✅ Aliyun SMS 简单处理
- ✅ 日志分片(异步日志文件管理)，
- ✅ 时间工具
- ✅ 字符串工具
- ✅ 数组工具
- ✅ 系统信息
- ✅ 超高精度Float64加、减、乘、除运算
- ✅ Bit  Bytes 单位换算工具
- ⚠️ file upload 文件上传
- ❌ error
- ❗ warning tip
# 6. 感谢
`gtbox`收到了 [Jetbrains](https://jb.gg/OpenSourceSupport) Licenses，非常感谢`Jetbrains`对开源项目的支持。

![JetBrains logo](./resource/img/jetBrains.svg)