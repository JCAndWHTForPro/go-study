// 第 2 课：变量、常量与基本类型
// 运行方式：go run ./02-variables
package main

import "fmt"

// 包级别常量：const 定义常量，值在编译期确定
// Pi 没写类型 —— 这是「无类型常量」，用到哪里就临时适配成哪个类型
const Pi = 3.14159

// MaxAge 显式写了 int —— 这是「有类型常量」，类型被钉死成 int
const MaxAge int = 100

// iota：自增常量生成器，常用于做“枚举”
// 在一个 const 块中，iota 从 0 开始，每行自动 +1
const (
	Sunday    = iota // 0
	Monday           // 1
	Tuesday          // 2
	Wednesday        // 3
)

func main() {
	// ===== 1. 四种变量声明方式 =====
	var name string = "Tom" // 完整写法：var 变量名 类型 = 值
	var age = 18            // 省略类型，由右值自动推断为 int
	var height float64      // 只声明不赋值，得到“零值” 0
	score := 95             // 短声明 :=（最常用，只能在函数内用）
	/*
		1、var 包内、函数内都能用；:= 只能在函数内
		2、:= 左边至少要有一个是新变量，否则报错
		3、只声明不赋值时，一定有零值（不会是"未定义"）

	*/

	fmt.Println("name:", name, "age:", age, "height:", height, "score:", score)

	// ===== 1.1 批量声明的几种写法 =====
	// 方式 A：一行声明多个同类型变量
	var x, y, z int = 1, 2, 3

	// 方式 B：一行声明多个不同类型（类型推断，省略类型）
	var p, q = 10, "hello"

	// 方式 C：短声明批量（函数内最常用）
	one, two := 1, "two"

	fmt.Println("批量声明:", x, y, z, p, q, one, two)

	// ===== 1.2 多返回值与 _ 丢弃值 =====
	// Go 函数常返回多个值，用 := 一次性接收
	// _ 是“空白标识符”，用来丢弃不需要的返回值（拿到的值若不用会编译报错）
	length, ok := lookup("go")
	fmt.Println("查到长度:", length, "是否命中:", ok)

	_, onlyOk := lookup("rust") // 只关心是否命中，长度用 _ 丢弃
	fmt.Println("只看是否命中:", onlyOk)

	// ===== 2. 零值演示：Go 没有“未初始化”，一切都有默认值 =====
	var (
		i int     // 0
		f float64 // 0
		b bool    // false
		s string  // "" 空字符串
	)
	fmt.Printf("零值 -> int:%d float:%.1f bool:%t string:%q\n", i, f, b, s)

	// ===== 3. 常量与 iota =====
	fmt.Println("Pi =", Pi)
	fmt.Println("周几枚举:", Sunday, Monday, Tuesday, Wednesday)

	// ===== 3.1 有类型常量 vs 无类型常量 =====
	// 无类型常量 Pi：用到哪里就临时适配成哪个类型，非常灵活
	var pi32 float32 = Pi // Pi 在这里当 float32 用
	var pi64 float64 = Pi // Pi 在这里当 float64 用
	fmt.Printf("无类型常量 Pi 适配: float32=%v float64=%v\n", pi32, pi64)

	// 有类型常量 MaxAge：类型被钉死成 int，规则和变量一样严格
	// var bigAge int64 = MaxAge       // ❌ 报错：int 不能直接给 int64
	var bigAge int64 = int64(MaxAge) // ✅ 必须显式转换
	fmt.Println("有类型常量 MaxAge 需显式转换:", bigAge)

	// ===== 4. 类型必须显式转换（Go 不做隐式转换）=====
	var a int = 10
	var c int64 = 20
	// sum := a + c        // ❌ 这行会编译报错：int 和 int64 不能直接相加
	sum := int64(a) + c // ✅ 必须显式把 a 转成 int64
	fmt.Println("显式转换后相加:", sum)

	// ===== 5. byte 与 rune =====
	var ch byte = 'A' // byte = uint8，存 ASCII
	var r rune = '中' // rune = int32，存一个 Unicode 字符
	fmt.Printf("byte 'A' 的值=%d, rune '中' 的码点=%d\n", ch, r)

	// ===== 6. 用 %T 查看变量的类型 =====
	fmt.Printf("age 的类型是: %T, score 的类型是: %T\n", age, score)
}

// lookup 演示「多返回值」：返回单词长度，以及是否查到（ok 模式）
// Go 里这种「值 + ok 布尔」的双返回值非常常见（如 map 取值、类型断言）
func lookup(word string) (int, bool) {
	dict := map[string]int{
		"go": 2,
		"hi": 2,
	}
	length, ok := dict[word] // map 取值天然返回两个值：值 + 是否存在
	return length, ok
}

// ============================================================
// 【补充】fmt 格式化打印函数族总结
// ============================================================
//
// 一、三大家族（按输出目标区分）
// ---------------------------------------------------------------
// 函数           | 输出到哪         | 返回值
// --------------|-----------------|-----------------------------
// fmt.Print     | 标准输出(终端)    | 写入字节数, error
// fmt.Sprint    | 返回 string      | string
// fmt.Fprint    | 写入 io.Writer   | 写入字节数, error
// ---------------------------------------------------------------
//
// 每个家族都有三个变体：
//   Print / Println / Printf    → 输出到终端
//   Sprint / Sprintln / Sprintf → 返回字符串
//   Fprint / Fprintln / Fprintf → 写入文件/网络等 io.Writer
//
// 助记：
//   S = String（结果存变量）
//   F = File/Writer（结果写到指定目标）
//   无前缀 = 直接打印到终端
//   ln 后缀 = 自动换行 + 空格分隔
//   f 后缀 = 支持格式化占位符
//
// 二、常用格式化占位符（Printf / Sprintf 的 %verb）
// ---------------------------------------------------------------
// 通用：
//   %v   默认格式（万能，啥类型都行）
//   %+v  结构体带字段名：{Name:Tom Age:18}
//   %#v  Go 语法格式：main.User{Name:"Tom", Age:18}
//   %T   打印变量的类型
//   %%   输出一个字面 % 号
//
// 整数：
//   %d   十进制
//   %b   二进制
//   %o   八进制
//   %x   十六进制（小写）
//   %X   十六进制（大写）
//
// 浮点数：
//   %f   小数形式（默认6位小数）
//   %.2f 保留2位小数
//   %e   科学计数法
//
// 字符串/字符：
//   %s   原样输出字符串
//   %q   带双引号的字符串（转义不可见字符）
//   %c   字符（rune → 对应的字符）
//
// 布尔：
//   %t   true / false
//
// 指针：
//   %p   指针地址
//
// 三、Errorf —— 特殊成员
// ---------------------------------------------------------------
// fmt.Errorf("余额不足: %d", balance) → 返回 error 类型
// 内部就是 Sprintf + 包装成 error，构造错误信息时非常常用
// 支持 %w 包装原始错误：fmt.Errorf("查询失败: %w", err)
// ============================================================
