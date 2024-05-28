package tools

import (
	"math/rand"
	"time"
)

var privateRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// 生成随机数，numberRange参数可填1-2个，1=[0~num)，2=[num1,num2)
func RanNumber(numberRange ...int) int {
	if len(numberRange) > 1 {
		return RanWithRange(0, numberRange[0]-1)
	} else {
		return RanWithRange(numberRange[0], numberRange[1])
	}
}

// 根据给定范围取随机数，结果为：min ~ max
func RanWithRange(min, max int) int {
	if min == max {
		return min
	} else if min > max {
		min, max = max, min
	}
	return min + privateRand.Intn(max-min+1)
}

// 根据最大值取随机数，结果为：0 ~ max-1，常用于数组随机下标 max=len(arr)
func RanWithMax(max int) int {
	if max == 0 {
		return 0
	}
	return privateRand.Intn(max)
}
