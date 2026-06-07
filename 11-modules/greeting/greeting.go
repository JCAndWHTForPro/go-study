// greeting 包：演示多包项目 + init 执行顺序
// 导入路径：go-study/11-modules/greeting
package greeting

import "fmt"

// init 函数：被依赖的包先 init
func init() {
	fmt.Println("[init] greeting 包的 init 被调用")
}

// Hello 返回问候语
func Hello(name string) string {
	return fmt.Sprintf("你好, %s! 这是来自 greeting 包的问候", name)
}
