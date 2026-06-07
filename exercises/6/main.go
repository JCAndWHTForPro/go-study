package main

import (
	"fmt"
	"math"
	"sort"
)

/*
## 第 6 课：接口

### 题 6.1 面积计算器（⭐⭐）

定义一个 `Shape` 接口，包含 `Area() float64` 方法。
实现三个类型：`Circle`、`Rectangle`、`Triangle`，各自实现 `Area()`。
写一个函数 `totalArea(shapes []Shape) float64` 计算总面积。

### 题 6.2 排序接口（⭐⭐⭐）

自定义一个 `StringSlice` 类型（底层是 `[]string`），实现 `sort.Interface` 的三个方法（`Len`、`Less`、`Swap`），让字符串按长度排序（短的在前）。
*/

type Shape interface {
	Area() float64
}
type Circle struct {
	r float64
}
type Rectangle struct {
	length float64
	witdh  float64
}
type Triangle struct {
	base   float64
	height float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.r * c.r
}
func (r Rectangle) Area() float64 {
	return r.witdh * r.length
}
func (t Triangle) Area() float64 {
	return (t.base * t.height) * 0.5
}

type StringSlice []string

func (s StringSlice) Len() int {
	//TODO implement me
	return len(s)
}

func (s StringSlice) Less(i, j int) bool {
	//TODO implement me
	return len(s[i]) < len(s[j])
}

func (s StringSlice) Swap(i, j int) {
	//TODO implement me
	s[i], s[j] = s[j], s[i]
}

func main() {
	ss := StringSlice{"1121", "112", "1"}
	sort.Sort(ss)
	fmt.Println(ss)
}

func totalArea(shapes []Shape) float64 {
	var sum float64
	for _, shape := range shapes {
		sum += shape.Area()
	}
	return sum
}
