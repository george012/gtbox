package gtbox_array

import (
	"fmt"
	"sync"
)

// Array 可变长数组
type Array struct {
	array []int      // 固定大小的数组，用满容量和满大小的切片来代替
	len   int        // 真正长度
	cap   int        // 容量
	lock  sync.Mutex // 为了并发安全使用的锁
}

// MakeArray 新建一个可变长数组
// 时间复杂度为：O(1)，因为分配内存空间和设置几个值是常数时间。
func MakeArray(len, cap int) *Array {
	s := new(Array)
	if len > cap {
		fmt.Printf("数组的长度不能大于容量")
	}

	// 把切片当数组用
	array := make([]int, cap, cap)

	// 元数据
	s.array = array
	s.cap = cap
	s.len = 0
	return s
}

// Append 增加一个元素
// 添加元素中，耗时主要在老数组中的数据移动到新数组，时间复杂度为：O(n)
func (a *Array) Append(element int) {
	// 并发锁
	a.lock.Lock()
	defer a.lock.Unlock()

	// 大小等于容量，表示没多余位置了
	if a.len == a.cap {
		// 没容量，数组要扩容，扩容到两倍
		newCap := 2 * a.len

		// 如果之前的容量为0，那么新容量为1
		if a.cap == 0 {
			newCap = 1
		}

		newArray := make([]int, newCap, newCap)

		// 把老数组的数据移动到新数组
		for k, v := range a.array {
			newArray[k] = v
		}

		// 替换数组
		a.array = newArray
		a.cap = newCap

	}

	// 把元素放在数组里
	a.array[a.len] = element
	// 真实长度+1
	a.len = a.len + 1

}

// AppendMany 增加多个元素
func (a *Array) AppendMany(element ...int) {
	for _, v := range element {
		a.Append(v)
	}
}

// Get 获取某个下标的元素
// 因为只获取下标的值，所以时间复杂度为 O(1)
func (a *Array) Get(index int) int {
	// 越界了
	if a.len == 0 || index >= a.len {
		fmt.Printf("索引超过了长度")
	}
	return a.array[index]
}

// Len 返回真实长度
// 时间复杂度为 O(1)
func (a *Array) Len() int {
	return a.len
}

// Cap 返回容量
// 时间复杂度为 O(1)
func (a *Array) Cap() int {
	return a.cap
}

// ToString 转换为字符串输出，主要用于打印
func (a *Array) ToString() (result string) {
	result = "["
	for i := 0; i < a.Len(); i++ {
		// 第一个元素
		if i == 0 {
			result = fmt.Sprintf("%s%d", result, a.Get(i))
			continue
		}

		result = fmt.Sprintf("%s %d", result, a.Get(i))
	}
	result = result + "]"
	return
}

// GTArrayReverse 数组倒序
func GTArrayReverse(arr *[]string) {
	length := len(*arr)
	var temp string
	for i := 0; i < length/2; i++ {
		temp = (*arr)[i]
		(*arr)[i] = (*arr)[length-1-i]
		(*arr)[length-1-i] = temp
	}
}

// IsContainWithStrArray 判断 字符串数组 是否包含 某个字符串
func IsContainWithStrArray(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
