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