/*
Package gtbox_color en: Color Tools, zh-cn: 颜色工具库
*/
package gtbox_color

// ANSIColor en: ANSI Color Enums, zh-cn: ANSI颜色枚举
type ANSIColor string

const (
	// Foreground Colors 前景色
	ANSIColorForegroundBlack   ANSIColor = "\033[30m" // en: Black, zh-cn: 前景色->黑色
	ANSIColorForegroundRed     ANSIColor = "\033[31m" // en: Red, zh-cn: 前景色->红色
	ANSIColorForegroundGreen   ANSIColor = "\033[32m" // en: Green, zh-cn: 前景色->绿色
	ANSIColorForegroundYellow  ANSIColor = "\033[33m" // en: Yellow, zh-cn: 前景色->黄色
	ANSIColorForegroundBlue    ANSIColor = "\033[34m" // en: Blue, zh-cn: 前景色->蓝色
	ANSIColorForegroundMagenta ANSIColor = "\033[35m" // en: Magenta, zh-cn: 前景色->洋红色
	ANSIColorForegroundCyan    ANSIColor = "\033[36m" // en: Cyan, zh-cn: 前景色->青色
	ANSIColorForegroundWhite   ANSIColor = "\033[37m" // en: White, zh-cn: 前景色->白色

	// Bright Foreground Colors 亮色前景色
	ANSIColorForegroundBrightBlack   ANSIColor = "\033[90m" // en: Bright Black, zh-cn: 亮色前景色->亮黑色
	ANSIColorForegroundBrightRed     ANSIColor = "\033[91m" // en: Bright Red, zh-cn: 亮色前景色->亮红色
	ANSIColorForegroundBrightGreen   ANSIColor = "\033[92m" // en: Bright Green, zh-cn: 亮色前景色->亮绿色
	ANSIColorForegroundBrightYellow  ANSIColor = "\033[93m" // en: Bright Yellow, zh-cn: 亮色前景色->亮黄色
	ANSIColorForegroundBrightBlue    ANSIColor = "\033[94m" // en: Bright Blue, zh-cn: 亮色前景色->亮蓝色
	ANSIColorForegroundBrightMagenta ANSIColor = "\033[95m" // en: Bright Magenta, zh-cn: 亮色前景色->亮紫色(品红)
	ANSIColorForegroundBrightCyan    ANSIColor = "\033[96m" // en: Bright Cyan, zh-cn: 亮色前景色->亮青色
	ANSIColorForegroundBrightWhite   ANSIColor = "\033[97m" // en: Bright White, zh-cn: 亮色前景色->亮白色

	// Background Colors 背景色
	ANSIColorBackgroundBlack   ANSIColor = "\033[40m" // en: Black, zh-cn: 背景色->黑色
	ANSIColorBackgroundRed     ANSIColor = "\033[41m" // en: Red, zh-cn: 背景色->红色
	ANSIColorBackgroundGreen   ANSIColor = "\033[42m" // en: Green, zh-cn: 背景色->绿色
	ANSIColorBackgroundYellow  ANSIColor = "\033[43m" // en: Yellow, zh-cn: 背景色->黄色
	ANSIColorBackgroundBlue    ANSIColor = "\033[44m" // en: Blue, zh-cn: 背景色->蓝色
	ANSIColorBackgroundMagenta ANSIColor = "\033[45m" // en: Magenta, zh-cn: 背景色->洋红色
	ANSIColorBackgroundCyan    ANSIColor = "\033[46m" // en: Cyan, zh-cn: 背景色->青色
	ANSIColorBackgroundWhite   ANSIColor = "\033[47m" // en: White, zh-cn: 背景色->白色

	// Bright Background Colors 亮色背景色
	ANSIColorBackgroundBrightBlack   ANSIColor = "\033[100m" // en: Bright Black, zh-cn: 亮色背景色->亮黑色
	ANSIColorBackgroundBrightRed     ANSIColor = "\033[101m" // en: Bright Red, zh-cn: 亮色背景色->亮红色
	ANSIColorBackgroundBrightGreen   ANSIColor = "\033[102m" // en: Bright Green, zh-cn: 亮色背景色->亮绿色
	ANSIColorBackgroundBrightYellow  ANSIColor = "\033[103m" // en: Bright Yellow, zh-cn: 亮色背景色->亮黄色
	ANSIColorBackgroundBrightBlue    ANSIColor = "\033[104m" // en: Bright Blue, zh-cn: 亮色背景色->亮蓝色
	ANSIColorBackgroundBrightMagenta ANSIColor = "\033[105m" // en: Bright Magenta, zh-cn: 亮色背景色->亮紫色(品红)
	ANSIColorBackgroundBrightCyan    ANSIColor = "\033[106m" // en: Bright Cyan, zh-cn: 亮色背景色->亮青色
	ANSIColorBackgroundBrightWhite   ANSIColor = "\033[107m" // en: Bright White, zh-cn: 亮色背景色->亮白色

	// Reset Color 重置颜色
	ANSIColorReset ANSIColor = "\033[0m" // en: Reset, zh-cn: 重置颜色设置
)
