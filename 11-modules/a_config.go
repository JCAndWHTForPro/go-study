// a_config.go：演示同包多文件 init 执行顺序
// 文件名以 a_ 开头，字母序最小 → 它的 init 最先执行
package main

import "fmt"

func init() {
	fmt.Println("[init] a_config.go 的 init（文件名字母序最小，最先执行）")
}

// 同一个文件可以有多个 init，按从上到下顺序执行
func init() {
	fmt.Println("[init] a_config.go 的第 2 个 init（同文件从上到下）")
}
