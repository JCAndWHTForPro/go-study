package request

/*

把第 18~19 题的博客系统改成 Gin API 服务：

- `GET /api/users/:id/posts` — 查某用户的所有文章（Preload 评论）
- `POST /api/posts` — 发文章（JSON: title, content, user_id）
- `POST /api/posts/:id/comments` — 给文章评论
- 加一个日志中间件记录请求耗时
- 路由分组 `/api`
*/
