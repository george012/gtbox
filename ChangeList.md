gtboX is a develop tools
<!-- TOC -->

- [1. v0.1.13](#1-v0113)
- [2. v0.0.73](#2-v0073)
- [3. v0.0.65](#3-v0065)
- [4. v0.0.52](#4-v0052)
- [5. v0.0.49](#5-v0049)
- [6. v0.0.48](#6-v0048)
- [7. v0.0.47](#7-v0047)
- [8. v0.1.54](#8-v0154)
- [9. v0.1.55](#9-v0155)
- [10. v0.1.56](#10-v0156)
- [11. v0.1.61](#11-v0161)
- [12. v0.1.62](#12-v0162)
- [13. v0.1.69](#13-v0169)

<!-- /TOC -->
# 1. v0.1.13
* 修复 `cmd` 工具windows执行错误

# 2. v0.0.73
* 修复 `加密库` 工具

# 3. v0.0.65
* 新增 `fyne` GUI工具

# 4. v0.0.52
* 新增 `LogF` 快捷含自定义Module封装的日志打印方法

# 5. v0.0.49
* 新增 `decimal` 和 `big.float` 转换
* 新增 `decimal` 和 `float64` 转换

# 6. v0.0.48
* 新增  `func GTBigFLoat2Float64(bigFloat *big.Float) float64{}` 方法 在 `gtbox_hashrate`

# 7. v0.0.47
* 增加 `ChangeLog.md` 即 更新日志记录文件
* 修改 `gtbox_hashrate` 传入`hs` 值为 `*big.FLoat`类型

# 8. v0.1.54
* 扩展 `gtbox_log` 默认单实，例扩展支持多实例，方便多扩快日志分离
* 更新依赖版本
* 扩展encryption并发安全锁

# 9. v0.1.55
* 修复 `gtbox_log` 默认日志路径处理逻辑

# 10. v0.1.56
* 扩展 `gtbox_log` 日志分片支持文件夹归集
* 新增 `gtbox_app` App通用属性，支持继承扩展

# 11. v0.1.61
* Fix `gtbox_log` 日志分片支持文件夹按日期归集的Bug
* 
# 12. v0.1.62
* Fix `gtbox_log` 修复日志路径拼接错误
# 13. v0.1.69
* Fix 新增`gtbox_redis`