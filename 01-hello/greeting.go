// 演示「多文件同包」：本文件和 main.go 在同一目录、同属 package main
//
// 要点：
// 1) 文件名叫 greeting.go（和 main.go 不同），但 package 仍必须是 main
// 2) 同包内的函数可以「直接调用」，不需要 import，也不用加包名前缀
// 3) 编译时同目录下所有 .go 文件会被当成一个整体
package main

import "fmt"

// greet 根据名字拼出一句问候语并返回
// 首字母小写 greet：包内私有，只能在本包（本目录）内使用
func greet(name string) string {
	return fmt.Sprintf("你好, %s! 欢迎学习 Go", name)
}
