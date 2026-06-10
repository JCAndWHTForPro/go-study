package main

import (
	"go-study/exercises/18-19/db"
	"go-study/exercises/18-19/request"
)

/*
## 第 18-19 课：GORM

### 题 18-19.1 博客系统模型（⭐⭐⭐）

用 GORM 实现一个简易博客的数据模型：

  - `User`（id, username, email）
  - `Post`（id, title, content, user_id）→ 属于 User
  - `Comment`（id, content, post_id, user_id）→ 属于 Post 和 User
    实现：

1. AutoMigrate 建表
2. 创建 2 个用户，每人写 2 篇文章，每篇文章 2 条评论
3. Preload 查询：查某个用户的所有文章及每篇文章的评论
4. Joins 查询：查某篇文章的所有评论及评论者名字

---

## 第 19 课：Gin

### 题 19.1 Gin + GORM 整合（⭐⭐⭐）

把第 18-19 题的博客系统改成 Gin API 服务：

- `GET /api/users/:id/posts` — 查某用户的所有文章（Preload 评论）
- `POST /api/posts` — 发文章（JSON: title, content, user_id）
- `POST /api/posts/:id/comments` — 给文章评论
- 加一个日志中间件记录请求耗时
- 路由分组 `/api`
*/
func main() {
	db.InitDB("root:root123@tcp(127.0.0.1:3306)/go_study?charset=utf8mb4&parseTime=True")
	request.InitRequestRegister()
}
