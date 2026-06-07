package main

import (
	"fmt"
	"regexp"
	"sort"
)

/*
## 第 5 课：复合类型

### 题 5.1 切片去重（⭐⭐）

写一个函数 `unique(nums []int) []int`，返回去重后的切片（保持原顺序）。
示例：`unique([]int{1,3,2,3,1,4})` → `[1,3,2,4]`

### 题 5.2 单词频率统计（⭐⭐）

给定一个字符串 `"go is great go is fast go go go"`，用 map 统计每个单词出现的次数，按频率从高到低排序输出。

### 题 5.3 矩阵转置（⭐⭐⭐）

写一个函数 `transpose(matrix [][]int) [][]int`，将 M×N 矩阵转置为 N×M。
示例：`[[1,2,3],[4,5,6]]` → `[[1,4],[2,5],[3,6]]`
*/
func main() {
	/*ints := unique([]int{1, 3, 2, 3, 1, 4})
	fmt.Println("结果是", ints)*/
	//countWord("go is great go is fast go go go")
	fmt.Println("结果", transpose([][]int{{1, 2, 3}, {4, 5, 6}}))
}

func unique(nums []int) []int {
	resultNums := []int{}
	contains := map[int]bool{}
	for _, num := range nums {
		if contains[num] {
			continue
		}
		contains[num] = true
		resultNums = append(resultNums, num)
	}
	return resultNums
}
func countWord(words string) {
	compile := regexp.MustCompile(`\s+`)
	wordSplit := compile.Split(words, -1)
	countMap := map[string]int{}
	for _, str := range wordSplit {
		countMap[str]++
	}
	keys := make([]string, 0, len(countMap))
	for k := range countMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return countMap[keys[i]] > countMap[keys[j]]
	})
	for _, k := range keys {
		fmt.Printf("%s: %d\n", k, countMap[k])
	}
}

func transpose(matrix [][]int) [][]int {
	if matrix == nil || len(matrix) == 0 || len(matrix[0]) == 0 {
		return matrix
	}
	m, n := len(matrix), len(matrix[0])
	result := make([][]int, n)
	for i := range result {
		result[i] = make([]int, m)
	}
	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
			result[j][i] = matrix[i][j]
		}
	}
	return result
}
