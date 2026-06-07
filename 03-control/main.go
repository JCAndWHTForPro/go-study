// 第 3 课：流程控制（if / for / switch）
// 运行方式：go run ./03-control
package main

import "fmt"

func main() {
	// ===== 1. if 条件判断 =====
	// 要点：条件不加括号 ()，但 { } 不能省略
	score := 85
	if score >= 90 {
		fmt.Println("优秀")
	} else if score >= 60 {
		fmt.Println("及格")
	} else {
		fmt.Println("不及格")
	}

	// ===== 1.1 if 的「初始化语句」=====
	// if 可以先执行一个短语句，再判断；变量作用域仅限于这个 if
	// 这种写法在「错误处理」里极其常见：if v, err := f(); err != nil { ... }
	if half := score / 2; half > 40 {
		fmt.Println("一半分数大于 40, half =", half)
	}
	// 这里访问不到 half，它只活在上面的 if 块里

	// ===== 2. for 循环（Go 唯一的循环关键字）=====
	// Go 没有 while，所有循环都用 for

	// 2.1 经典三段式：初始化; 条件; 后置
	sum := 0
	for i := 1; i <= 5; i++ {
		sum += i
	}
	fmt.Println("1 到 5 求和 =", sum)

	// 2.2 当 while 用：只保留条件
	n := 3
	for n > 0 {
		fmt.Println("倒计时:", n)
		n--
	}

	// 2.3 死循环：用 break 跳出
	count := 0
	for {
		count++
		if count >= 3 {
			break // 跳出整个循环
		}
	}
	fmt.Println("死循环跑了", count, "次后 break")

	// 2.4 for-range 遍历：遍历切片、map、字符串等
	fruits := []string{"apple", "banana", "cherry"}
	for index, value := range fruits {
		fmt.Printf("第 %d 个水果是 %s\n", index, value)
	}
	// 只要值不要下标，用 _ 丢弃下标
	for _, value := range fruits {
		fmt.Println("只要值:", value)
	}

	// continue：跳过本次，进入下一次循环
	for i := 1; i <= 5; i++ {
		if i%2 == 0 {
			continue // 跳过偶数
		}
		fmt.Println("奇数:", i)
	}

	// ===== 3. switch 分支 =====
	// 要点：Go 的 case 默认「自动 break」，不会像 C/Java 那样穿透
	day := 3
	switch day {
	case 1, 2, 3, 4, 5: // 一个 case 可以匹配多个值
		fmt.Println("工作日")
	case 6, 7:
		fmt.Println("周末")
	default:
		fmt.Println("非法的星期")
	}

	// 3.1 无表达式的 switch：当「多分支 if-else」用，更清晰
	temp := 28
	switch {
	case temp >= 30:
		fmt.Println("天气炎热")
	case temp >= 20:
		fmt.Println("天气舒适")
	default:
		fmt.Println("天气偏凉")
	}

	// 3.2 fallthrough：手动让 case 穿透到下一个（少用）
	switch num := 1; num {
	case 1:
		fmt.Println("匹配到 1")
		fallthrough // 强制继续执行下一个 case，不再判断条件
	case 2:
		fmt.Println("被 fallthrough 带到了 2")
	case 3:
		fmt.Println("这一行不会执行")
	}
}
