// calc 包：演示多包项目中的子包
// 导入路径：go-study/11-modules/calc
package calc

import "fmt"

// init 函数：在 main() 之前自动执行
func init() {
	fmt.Println("[init] calc 包的 init 被调用")
}

// Add 导出函数（首字母大写）
func Add(a, b int) int {
	return a + b
}

// Sub 导出函数
func Sub(a, b int) int {
	return a - b
}

// multiply 未导出函数（首字母小写），只能在 calc 包内使用
func multiply(a, b int) int {
	return a * b
}
