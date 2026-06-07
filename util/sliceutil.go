package util

// Make2DIntSlice 创建一个 rows×cols 的二维 int 切片，所有元素初始化为 0
func Make2DIntSlice(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}
	return matrix
}
