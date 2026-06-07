package main

import "fmt"

/*
### 题 1.1 温度转换器（⭐）

写一个程序，定义一个华氏温度变量 `fahrenheit = 98.6`，将它转为摄氏温度并打印。
公式：`celsius = (fahrenheit - 32) * 5 / 9`
要求：用 `:=` 短声明，打印时保留 2 位小数。

### 题 1.2 多返回值交换（⭐）

声明两个变量 `a = 10, b = 20`，不使用第三个变量，利用 Go 的多赋值特性一行完成交换，然后打印交换后的值。

### 题 1.3 iota 星期枚举（⭐⭐）

用 `const + iota` 定义一周七天的枚举（Sunday=0, Monday=1, ..., Saturday=6），
然后写一个函数 `dayName(day int) string`，输入数字返回中文星期名。测试打印 Monday 和 Friday 对应的中文。
*/

const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func main() {
	//PrintTemperatureToCelsius()
	//SwapABPrint()
	name, err := DayName(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(name)
}

func PrintTemperatureToCelsius() {
	fahrenheit := 98.6
	celsius := (fahrenheit - 32) * 5 / 9
	fmt.Printf("当前的摄氏温度是：%.2f", celsius)
}

func DayName(day int) (string, error) {
	switch day {
	case Sunday:
		return "星期日", nil
	case Monday:
		return "星期一", nil
	case Tuesday:
		return "星期二", nil
	case Wednesday:
		return "星期三", nil
	case Thursday:
		return "星期四", nil
	case Friday:
		return "星期五", nil
	case Saturday:
		return "星期六", nil
	}
	return "", fmt.Errorf("你当前的不是我们的正常星期数字：%d", day)
}

func SwapABPrint() {
	a, b := 10, 20
	a, b = b, a
	fmt.Printf("交换后的ab是：%d,%d", a, b)

}
