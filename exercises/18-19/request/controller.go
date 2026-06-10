package request

import (
	"fmt"
	"go-study/exercises/18-19/db"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

/*

把第 18-19 题的博客系统改成 Gin API 服务：

- `GET /api/users/:id/posts` — 查某用户的所有文章（Preload 评论）
- `POST /api/posts` — 发文章（JSON: title, content, user_id）
- `POST /api/posts/:id/comments` — 给文章评论
- 加一个日志中间件记录请求耗时
- 路由分组 `/api`
*/

func InitRequestRegister() error {
	engine := gin.Default()

	group := engine.Group("/api")
	group.Use(logMiddleware)
	{
		group.GET("/users/:id/posts", userPostGet)
		group.POST("/posts", createPost)
		group.POST("/posts/:id/comments", createComment)

	}

	err := engine.Run(":9090")

	if err != nil {
		return fmt.Errorf("异常：%w", err)
	}
	return nil
}

func createComment(context *gin.Context) {
	postId, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": "获取最新的id失败，请输入数字类型的id",
		})
		return
	}
	var comment db.Comment
	if err = context.ShouldBindJSON(&comment); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("参数异常:%v", err),
		})
		return
	}
	comment.PostID = uint(postId)
	u, err := db.CreateComment(&comment)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("异常创建评论:%v", err),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"msg":  "创建评论成功",
		"data": u,
	})

}

func createPost(context *gin.Context) {
	var post db.Post
	if err := context.ShouldBindJSON(&post); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("参数异常:%v", err),
		})
		return
	}
	postId, err := db.CreatePost(&post)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("创建异常:%v", err),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"msg":  "创建成功",
		"data": postId,
	})
}

func logMiddleware(context *gin.Context) {
	now := time.Now()
	context.Next()
	duration := time.Since(now)
	log.Printf("[自定义日志] %s %s → %d (%v)\n",
		context.Request.Method, context.Request.URL.Path, context.Writer.Status(), duration)

}

func userPostGet(context *gin.Context) {
	userId, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		fmt.Println("id获取异常，请传入正常的数字类型", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"msg": "id获取异常，请传入正常的数字类型",
		})
		return
	}
	comment := db.QueryAllUserPostAndComment(uint(userId))
	context.JSON(http.StatusOK, comment)

}
