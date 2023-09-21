<!-- TOC -->

- [1. 使用说明](#1-使用说明)
    - [1.1. 针对windows](#11-针对windows)
- [2. 安装、升级、卸载](#2-安装升级卸载)
- [3. 依赖](#3-依赖)
- [4. 二开Build](#4-二开build)
- [5. 必要支持](#5-必要支持)
- [6. 功能](#6-功能)
- [7. 删除所有本地和远端 tag](#7-删除所有本地和远端-tag)

<!-- /TOC -->

# 1. 使用说明
* <font color=red>只保留2个可运行版本，方便维护</font>
* 必须安装了`git`
* 项目必须用`go mod`自动管理依赖
* 必须：`go version` `>=` `1.18`

## 1.1. 针对windows
* Windows 开启administrator权限
* IDE 全局用户安装
* 如果目前已经是单用户安装 在`IDE` →`属性` → `兼容性`中设置`以管理员的身份启用`
* 内置 git config 强制 `LF`
* 强制 `LF` 设置
    ```
    git config --global core.autocrlf input
    
    git config --global core.safecrlf true
    ```

# 2. 安装、升级、卸载
* <font color=red>在任意`golang`项目根目录下使用`terminal`执行如下命令</font>
```
wget --no-check-certificate https://raw.githubusercontent.com/george012/gtbox/master/install_gtbox.sh && chmod a+x ./install_gtbox.sh && ./install_gtbox.sh
```

# 3. 依赖
```
// scp
go get -u github.com/bramvdbogaerde/go-scp@latest

// excel
go get -u github.com/qax-os/excelize/v2@latest
// req 包
go get -u github.com/imroc/req/v3@latest

// gjson
go get -u github.com/tidwall/gjson@latest

// ants
go get -u github.com/panjf2000/ants/v2@latest
// GBK和UTF-8转换
go get -u github.com/axgle/mahonia@latest
```

# 4. 二开Build
```
./build 
```
*   自动化打包、提交、打Tag、并删除提交冗余Tags

# 5. 必要支持
*   CGO支持
*   MAC安装最新版本Xcode及Command Line Tools


# 6. 功能
- [x] CGO支持
- [x] 自定义加、解密
- [x] 简单的 SSH Client
- [x] 简易 SCP 工具
- [x] 简单的 HTTP Client
- [x] 简单的 ORM 封装
- [x] Aliyun SMS 简单处理
- [x] 日志分片
- [x] 时间工具
- [x] 字符串工具
- [x] 数组工具
- [x] 系统信息
- [x] 超高精度Float64加、减、乘、除运算
- [x] Bit  Bytes 单位换算工具
- [ ] 跨平台GUI工具---Fyne
- [ ] 跨平台GUI工具---Wails


# 7. 删除所有本地和远端 tag
```
git push origin --delete $(git tag -l) && git tag -d $(git tag -l)
```