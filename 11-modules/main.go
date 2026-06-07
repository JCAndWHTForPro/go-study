// 第 11 课：包管理与项目组织（go mod / 多包 / init / internal）
// 运行方式：go run ./11-modules
//
// ============================================================
// 【知识点总结】Go 项目组织
// ------------------------------------------------------------
//
// 1) go.mod：项目的「身份证」
//    - module go-study        ← 模块名，import 路径的前缀
//    - go 1.21                ← 最低 Go 版本
//    - require xxx v1.2.3     ← 依赖列表（go get 自动添加）
//
// 2) 常用命令：
//    go mod init <模块名>    创建 go.mod
//    go mod tidy             自动整理依赖（加缺的、删多的）
//    go get <包路径>          添加/更新依赖
//    go get <包路径>@v1.2.3  指定版本
//    go mod vendor           把依赖复制到 vendor/ 目录（离线构建）
//
// 3) 目录结构惯例（Go 社区标准）：
//    myproject/
//    ├── go.mod              模块定义
//    ├── main.go             入口（package main）
//    ├── internal/           私有包（只能被本模块 import，外部不可用）
//    │   └── config/
//    │       └── config.go
//    ├── pkg/                公共库（可被外部 import）
//    │   └── util/
//    │       └── util.go
//    └── cmd/                多个可执行程序
//        ├── server/
//        │   └── main.go
//        └── cli/
//            └── main.go
//
// 4) init 函数：
//    - 每个包可以有一个或多个 init() 函数
//    - 在 main() 之前自动执行，不能手动调用
//    - 常用于：初始化配置、注册驱动、校验环境
//
//    完整执行顺序规则（三层）：
//    ┌──────────────────────────────────────────────────────────┐
//    │ 第 1 层：包之间 → 被依赖的包先 init（依赖树从叶子到根）    │
//    │ 第 2 层：同包多文件 → 按「文件名字母序」执行              │
//    │ 第 3 层：同文件多 init → 按「从上到下」出现顺序执行        │
//    │ 最后：所有 init 完毕后，才执行 main()                    │
//    └──────────────────────────────────────────────────────────┘
//
//    示例：main 包有三个文件 config.go / handler.go / main.go
//      ① config.go 的 init()        ← c 字母序最小
//      ② handler.go 的 init()       ← h 第二
//      ③ main.go 的第 1 个 init()    ← m 第三，文件内从上到下
//      ④ main.go 的第 2 个 init()    ← 同文件第二个
//      ⑤ main()                     ← 最后
//
//    ⚠️ 注意：
//    - 不要依赖 init 顺序写业务逻辑（官方不100%保证同包文件顺序永不变）
//    - init 里保持简单，复杂逻辑放 main 显式调用更清晰
//    - 建议每个文件最多一个 init，多了可读性差
//
// 5) internal 目录：
//    - 放在 internal/ 下的包只能被「同一模块」内的代码 import
//    - 外部模块即使知道路径也无法导入（编译器强制限制）
//    - 用来放「不想暴露给外部的内部实现」
//
// 6) 父子目录的关系（重要！和 Java 不同）：
//    Go 里「子目录 = 独立的包」，和父目录没有任何从属/继承关系。
//
//    11-modules/              ← package main（父目录）
//    ├── main.go
//    ├── calc/                ← package calc（子目录，完全独立的包）
//    │   └── calc.go
//    └── greeting/            ← package greeting（子目录，完全独立的包）
//        └── greeting.go
//
//    关键认知：
//    - 子目录不会"继承"父目录的 package 名，它有自己的 package 声明
//    - 父目录有 main，子目录也可以有 main（不同目录 = 不同包 = 不冲突）
//    - 一个 module 可以有多个 main 函数，只要在不同目录
//    - 子目录的包要被父目录使用，必须通过 import（和导入第三方包一样）
//
//    和 Java 对比：
//    Java:  com.example.app 和 com.example.app.util 有层级关系
//    Go:    myapp/ 和 myapp/util/ 是「平级独立」的包，没有层级继承
//
//    唯一的限制：同一个目录下所有 .go 文件必须是同一个 package 名
//    （第 1 课就学过的铁律，子目录和父目录不算同一个目录）
// ============================================================
package main

import (
	"fmt"

	// 导入同模块下的子包：模块名(go-study) + 目录路径
	"go-study/11-modules/calc"
	"go-study/11-modules/greeting"
)

// init 函数：在 main() 之前自动执行
// 可以有多个 init，按文件中出现顺序执行
func init() {
	fmt.Println("[init] main 包的 init 被调用（在 main() 之前）")
}

func main() {
	fmt.Println("[main] 程序开始\n")

	// ===== 1. 多包调用演示 =====
	// calc 包：导出函数首字母大写
	fmt.Println("3 + 5 =", calc.Add(3, 5))
	fmt.Println("10 - 4 =", calc.Sub(10, 4))

	// greeting 包
	fmt.Println(greeting.Hello("凌还"))

	// ===== 2. init 执行顺序 =====
	// 观察输出：greeting 的 init → calc 的 init → main 的 init → main()
	// 被依赖的包先 init，最后才是 main 包

	// ===== 3. 包的可见性回顾 =====
	// calc.Add     ✅ 首字母大写，导出
	// calc.add     ❌ 首字母小写，包内私有（编译报错）
	// internal/    ❌ 外部模块无法 import（编译器限制）

	fmt.Println("\n[main] 程序结束")
}
