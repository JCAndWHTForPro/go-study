// 第 1 课：Hello World
// 运行方式：go run ./01-hello
//
// ===== Go 代码组织的基础知识 =====
//
// 1) package（包）：每个 .go 文件第一行有效代码必须是 package 声明
//    可执行程序的入口包必须叫 main，且必须有一个 main() 函数作为程序起点
//
// 2) import（导入）：引入要用到的标准库或第三方包
//    - 单个：import "fmt"
//    - 多个用「分组」写法（推荐）：
//        import (
//            "fmt"
//            "os"
//        )
//    - 导入了就必须用，未使用的 import 会直接编译报错
//
// 3) 缩进：Go 官方统一用「Tab」缩进（不是空格）
//    不用自己纠结对齐，保存时跑 gofmt / go fmt 会自动格式化
//
// 4) 大括号 {：左大括号必须跟在同一行末尾，不能另起一行
//    ✅ func main() {        ❌ func main()
//          ...                    {   // 这样写会编译报错
//       }
//
// 5) 语句分隔符：Go 语句结尾「不写分号 ;」
//    编译器会在每行末尾自动插入分号，所以一行写一条语句即可
//    （这也是为什么 { 不能换行——换行会被自动插入分号而出错）
//
// 6) 注释：// 单行注释   /* ... */ 块注释
package main

import (
	"fmt"

	// 跨包导入：写「完整路径」= 模块名(go-study) + 目录(util)
	// 路径在 go.mod 里的 module 名基础上拼出来，不是写包名
	"go-study/util"
)

func main() {
	fmt.Println("hello world!")

	// 直接调用同包另一个文件 greeting.go 里的 greet 函数
	// 注意：不需要 import，也不用写 greeting.greet，因为它们同属 package main
	message := greet("凌还")
	fmt.Println(message)

	// 跨包调用：用的时候只写「包名.函数名」= util.Xxx（不是写路径）
	// 只能调用首字母大写的「导出」函数
	fmt.Println("反转后:", util.Reverse("你好Go"))
	fmt.Println("转大写:", util.ToUpper("hello"))
}
