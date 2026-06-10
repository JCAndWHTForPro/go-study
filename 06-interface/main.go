// 第 6 课：方法与接口（method / interface / 隐式实现）
// 运行方式：go run ./06-interface
//
// Go 没有 class、没有继承，面向对象靠「结构体 + 方法 + 接口」三件套实现。
// 和 Java 最大的不同：接口是「隐式实现」——不需要 implements 关键字，
// 只要一个类型实现了接口里所有方法，它就「自动」算实现了该接口（鸭子类型）。
//
// ============================================================
// 【总结】接口赋值 & any 参数的 值/指针 规则
// ------------------------------------------------------------
//
// ▶ 接口赋值：取决于方法是值接收者还是指针接收者实现的
//
//	值接收者实现：值和指针都能赋给接口
//	  func (c Circle) Area() float64 { ... }
//	  var s Shape = Circle{}   // ✅ 值 → 接口
//	  var s Shape = &Circle{}  // ✅ 指针 → 接口（值接收者的方法包含在指针的方法集里）
//
//	指针接收者实现：只有指针能赋给接口
//	  func (d *Dog) Speak() string { ... }
//	  var a Animal = &Dog{}    // ✅ 指针 → 接口
//	  var a Animal = Dog{}     // ❌ 编译报错！值的方法集里没有指针接收者的方法
//
//	口诀：值接收者→都行；指针接收者→只能指针。
//
// ▶ any（interface{}）参数：啥都能接，值和指针都行
//
//	func Decode(v any) error   // any 参数
//	Decode(user)    // ✅ 传值
//	Decode(&user)   // ✅ 传指针
//	两种都能编译通过，但传哪个取决于「业务需要」：
//	→ 如果函数需要修改你的变量（如 json.Decode、json.Unmarshal），必须传 &
//	  因为传值只是副本，函数改了也白改，传指针才能改到原变量
//	→ 如果函数只读不写（如 json.Marshal、fmt.Println），传值即可
//
//	常见的「必须传 &」的 any 参数函数：
//	  json.Unmarshal(data, &result)   // 要把 JSON 填进 result
//	  json.NewDecoder(r).Decode(&obj) // 要把 JSON 填进 obj
//	  fmt.Scan(&input)                // 要把用户输入写进 input
//	  db.QueryRow(...).Scan(&col)     // 要把查询结果写进 col
//
// ============================================================
package main

import "fmt"

// ===== 1. 方法：给结构体（或任意自定义类型）绑定函数 =====
type Rectangle struct {
	Width  float64
	Height float64
}

// 方法 = 普通函数 + 一个「接收者」(receiver)
// (r Rectangle) 就是接收者，表示这个方法属于 Rectangle 类型

// 值接收者：方法内拿到的是副本，改不动原对象，适合只读
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// 指针接收者：方法内能修改原对象，适合需要改字段的场景
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

// ===== 2. 接口：定义「一组方法」的集合，只关心「能做什么」=====
// Shape 接口：任何「有 Area() 和 Perimeter() 方法」的类型都算实现了它
type Shape interface {
	Area() float64
	Perimeter() float64
}

// 给 Rectangle 补上 Perimeter，让它满足 Shape 接口
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// 再定义一个 Circle，也实现 Shape 接口
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14159 * c.Radius
}

// describe 接收「接口」类型参数：传 Rectangle 还是 Circle 都行
// 这就是「多态」：同一个函数，传不同实现，表现不同
func describe(s Shape) {
	fmt.Printf("  面积=%.2f, 周长=%.2f\n", s.Area(), s.Perimeter())
}

func main() {
	// ===== 方法调用 =====
	rect := Rectangle{Width: 3, Height: 4}
	fmt.Println("矩形面积:", rect.Area())

	// 指针接收者方法：Go 会自动取地址，rect.Scale 等价于 (&rect).Scale
	rect.Scale(2)
	fmt.Println("放大 2 倍后:", rect, " 新面积:", rect.Area())

	// ===== 接口与多态 =====
	// 关键：Rectangle 和 Circle 都没写 implements，但都自动满足 Shape 接口
	shapes := []Shape{
		Rectangle{Width: 2, Height: 3},
		Circle{Radius: 5},
	}
	for i, s := range shapes {
		fmt.Printf("图形[%d]:\n", i)
		describe(s) // 多态：自动调用各自类型的 Area/Perimeter
	}

	// ===== 3. 空接口 any（interface{}）：可以存任意类型 =====
	// any 是 interface{} 的别名（Go 1.18-19+），相当于 Java 的 Object
	var anything any
	anything = 42
	fmt.Println("空接口存 int:", anything)
	anything = "hello"
	fmt.Println("空接口存 string:", anything)

	// ===== 4. 类型断言：把接口「还原」回具体类型 =====
	// 语法：value, ok := 接口变量.(目标类型)
	if str, ok := anything.(string); ok {
		fmt.Println("类型断言成功，是 string:", str)
	}

	// ===== 5. type switch：根据接口里的实际类型分支处理 =====
	printType(42)
	printType("go")
	printType(3.14)
	printType(true)

	// ===== 6. := 推断具体类型 vs var 显式接口类型 =====
	// := 短声明：左边不能标类型，类型由右边字面量「推断」出来
	c := Circle{Radius: 5} // c 的类型被推断为「具体类型 Circle」
	fmt.Printf(":= 推断 -> c 的类型是 %T\n", c)

	// var + 显式接口类型：把「具体类型」装进「接口变量」里
	// 想让变量是 Shape 接口类型，就必须用 var 显式标，:= 做不到
	var s Shape = Circle{Radius: 5} // s 的类型是「接口 Shape」
	fmt.Printf("var 显式 -> s 的类型是 %T\n", s)

	// 区别有什么用？变量是接口类型时，只能调接口里声明的方法，但能接收任意实现类
	// 比如 s 可以随时换成 Rectangle，c 不行（c 被钉死成 Circle）
	s = Rectangle{Width: 2, Height: 3} // ✅ 接口变量可换成另一个实现
	fmt.Printf("s 换成 Rectangle 后类型是 %T, 面积=%.2f\n", s, s.Area())
}

// printType 用 type switch 判断空接口里到底装的是什么类型
func printType(v any) {
	switch value := v.(type) {
	case int:
		fmt.Printf("这是 int: %d\n", value)
	case string:
		fmt.Printf("这是 string: %s\n", value)
	case float64:
		fmt.Printf("这是 float64: %.2f\n", value)
	default:
		fmt.Printf("未知类型: %v\n", value)
	}
}

// ============================================================
// 【补充】接口断言 & 类型转换 完整总结
// ============================================================
//
// 一、类型断言（Type Assertion）—— 从接口中提取具体类型
// ---------------------------------------------------------------
// 语法：value, ok := 接口变量.(目标类型)
//
// 1. 安全断言（推荐，带 ok 检查）
//    str, ok := v.(string)
//    if ok {
//        // 断言成功，str 是 string 类型
//    }
//
// 2. 非安全断言（不带 ok，失败直接 panic）
//    str := v.(string)  // 如果 v 不是 string → panic!
//
// 3. 断言到接口类型（判断是否实现了某个接口）
//    s, ok := v.(Shape)  // v 是否实现了 Shape 接口？
//    if ok {
//        fmt.Println(s.Area())  // 可以调 Shape 的方法了
//    }
//
// 二、type switch —— 多类型分支判断
// ---------------------------------------------------------------
// 语法：switch value := v.(type) { case 类型: ... }
//
//    switch value := v.(type) {
//    case int:          // value 自动是 int 类型
//    case string:       // value 自动是 string 类型
//    case Shape:        // value 自动是 Shape 接口类型
//    case nil:          // 接口变量是 nil
//    default:           // 都不匹配
//    }
//
// 注意：.(type) 只能用在 switch 里，不能单独写 v.(type)
//
// 三、类型断言 vs 类型转换 的区别
// ---------------------------------------------------------------
//
// 类型断言：接口 → 具体类型（从接口里"取出来"）
//   var v any = "hello"
//   s := v.(string)        // ✅ 接口 → string
//
// 类型转换：具体类型 → 具体类型（数值之间互转）
//   var a int = 42
//   b := float64(a)        // ✅ int → float64
//   c := int64(a)          // ✅ int → int64
//   s := string([]byte{})  // ✅ []byte → string
//
// 区别：
//   类型转换 float64(a)   → 左右都是具体类型，编译期确定
//   类型断言 v.(string)   → 左边必须是接口类型，运行时检查
//   两者语法不同，不能混用！
//
// 四、实际使用场景
// ---------------------------------------------------------------
//
// 场景 1：处理 any 参数
//   func process(v any) {
//       if num, ok := v.(int); ok { ... }
//   }
//
// 场景 2：errors.As（从 error 接口提取自定义错误类型）
//   var divErr *DivideError
//   if errors.As(err, &divErr) { ... }
//
// 场景 3：判断某个值是否实现了可选接口
//   if closer, ok := reader.(io.Closer); ok {
//       defer closer.Close()
//   }
//
// 场景 4：JSON 反序列化后的 any 值
//   var data any
//   json.Unmarshal([]byte(`{"name":"Tom"}`), &data)
//   m := data.(map[string]any)        // JSON 对象 → map
//   name := m["name"].(string)        // 取出具体值
// ============================================================
