/*
Package gtbox_array 基于线程安全的 可变长 Array/slice 封装
*/
package gtbox_array

import (
	"fmt"
	"sync"
)

// GTArray 基于线程安全的 可变长 Array/slice
type GTArray struct {
	array []interface{} // 能够存储任何类型的数组
	len   int64         // 真正长度
	cap   int64         // 容量
	lock  sync.RWMutex  // 读写锁，用于并发安全
}

// NewArray 新建一个可变长数组
func NewArray(len, cap int64) *GTArray {
	s := new(GTArray)
	if len > cap {
		fmt.Printf("数组的长度不能大于容量\n")
		return nil
	}
	array := make([]interface{}, cap, cap)
	s.array = array
	s.cap = cap
	s.len = len
	return s
}

// Append 增加一个元素
func (a *GTArray) Append(element interface{}) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.len == a.cap {
		newCap := 2 * a.len
		if a.cap == 0 {
			newCap = 1
		}
		newArray := make([]interface{}, newCap, newCap)
		copy(newArray, a.array)
		a.array = newArray
		a.cap = newCap
	}
	a.array[a.len] = element
	a.len = a.len + 1
}

// AppendMany 增加多个元素
func (a *GTArray) AppendMany(elements ...interface{}) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for _, v := range elements {
		a.Append(v)
	}
}

// Get 获取某个下标的元素
func (a *GTArray) Get(index int64) interface{} {
	a.lock.RLock()
	defer a.lock.RUnlock()
	if a.len == 0 || index >= a.len {
		fmt.Printf("索引超过了长度\n")
		return nil
	}
	return a.array[index]
}

// Len 返回真实长度
func (a *GTArray) Len() int64 {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.len
}

// Cap 返回容量
func (a *GTArray) Cap() int64 {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.cap
}

// ToString 转换为字符串输出，主要用于打印
func (a *GTArray) ToString() string {
	a.lock.RLock()
	defer a.lock.RUnlock()
	result := "["
	for i := int64(0); i < a.Len(); i++ {
		if i == 0 {
			result = fmt.Sprintf("%s%v", result, a.Get(i))
			continue
		}
		result = fmt.Sprintf("%s %v", result, a.Get(i))
	}
	result = result + "]"
	return result
}
