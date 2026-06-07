package main

import "fmt"

/*
## 第 10 课：泛型

### 题 10.1 泛型 Map/Filter（⭐⭐）

写两个泛型函数：

  - `Map[T, U](slice []T, fn func(T) U) []U` — 映射
  - `Filter[T](slice []T, fn func(T) bool) []T` — 过滤
    用例：把 `[]int{1,2,3,4,5}` 过滤出偶数，然后每个乘以 10。

### 题 10.2 泛型缓存（⭐⭐⭐）

实现一个泛型的简易缓存 `type Cache[K comparable, V any] struct`：

  - `Set(key K, value V)`
  - `Get(key K) (V, bool)`
  - `Delete(key K)`
    用 map 实现。
*/
type Cache[K comparable, V any] struct {
	CacheData map[K]V
}

func (cc *Cache[K, V]) Delete(key K) {
	delete(cc.CacheData, key)
}
func (cc *Cache[K, V]) Set(key K, value V) {
	cc.CacheData[key] = value
}
func (cc *Cache[K, V]) Get(key K) (V, bool) {
	vv, isExit := cc.CacheData[key]
	if !isExit {
		return vv, false
	}
	return vv, true
}
func main() {
	// === 题 10.1 测试 ===
	slice := []int{1, 2, 3, 4, 5}
	slice = Filter(slice, func(value int) bool {
		return value%2 == 0
	})
	slice = Map(slice, func(value int) int {
		return value * 10
	})
	fmt.Println("最后的结果是：", slice)

	// === 题 10.2 Cache 测试 ===
	cache := &Cache[string, int]{CacheData: make(map[string]int)}

	cache.Set("apple", 5)
	cache.Set("banana", 3)
	cache.Set("cherry", 8)
	fmt.Println("Set 三个水果后:", cache.CacheData)

	val, ok := cache.Get("banana")
	fmt.Printf("Get banana: value=%d, ok=%t\n", val, ok)

	_, ok = cache.Get("grape")
	fmt.Printf("Get grape (不存在): ok=%t\n", ok)

	cache.Delete("banana")
	fmt.Println("Delete banana 后:", cache.CacheData)

	_, ok = cache.Get("banana")
	fmt.Printf("再次 Get banana: ok=%t\n", ok)
}
func Filter[T any](slice []T, fn func(T) bool) []T {
	res := make([]T, 0, len(slice))
	for _, v := range slice {
		if !fn(v) {
			continue
		}
		res = append(res, v)
	}
	return res
}
func Map[T any, U any](slice []T, fn func(T) U) []U {
	result := []U{}
	for _, v := range slice {
		result = append(result, fn(v))
	}
	return result
}
