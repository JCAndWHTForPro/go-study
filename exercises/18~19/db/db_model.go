package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
用 GORM 实现一个简易博客的数据模型：

  - `User`（id, username, email）
  - `Post`（id, title, content, user_id）→ 属于 User
  - `Comment`（id, content, post_id, user_id）→ 属于 Post 和 User
    实现：

1. AutoMigrate 建表
2. 创建 2 个用户，每人写 2 篇文章，每篇文章 2 条评论
3. Preload 查询：查某个用户的所有文章及每篇文章的评论
4. Joins 查询：查某篇文章的所有评论及评论者名字
*/
type User struct {
	gorm.Model
	UserName string `gorm:"type:varchar(100);not null"`
	Email    string
}

type Post struct {
	gorm.Model
	Title    string
	Content  string
	UserID   uint
	User     User
	Comments []Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
	Post    Post
	UserID  uint
	User    User
}

var DB *gorm.DB

func queryPostCommentAdnCommentUser(postId uint) ([]Comment, map[uint]string) {
	var comments []Comment
	DB.Model(&Comment{}).Joins("Post").Joins("User").Where("post_id = ?", postId).Find(&comments)
	mapp := map[uint]string{}

	if len(comments) != 0 {
		for _, cc := range comments {
			mapp[cc.ID] = cc.User.UserName
		}
	}

	return comments, mapp

}
func queryUserPostAndComment(userId uint) ([]Post, map[uint][]Comment) {
	var posts []Post
	tx := DB.Where("user_id = ?", userId).Find(&posts)
	if tx.RowsAffected == 0 {
		return []Post{}, map[uint][]Comment{}
	}
	commentMap := map[uint][]Comment{}
	for _, p := range posts {
		var comments []Comment
		DB.Where("post_id = ?", p.ID).Find(&comments)
		commentMap[p.ID] = comments

	}
	return posts, commentMap
}
func createUser(user *User) uint {
	DB.Create(user)
	return user.ID
}

func initDB(dsn string) error {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库打开异常:%w", err)
	}
	err = DB.AutoMigrate(&User{})
	if err != nil {
		return fmt.Errorf("创建user表失败：%w", err)
	}
	err = DB.AutoMigrate(&Post{})
	if err != nil {
		return fmt.Errorf("创建post表失败：%w", err)
	}
	err = DB.AutoMigrate(&Comment{})
	if err != nil {
		return fmt.Errorf("创建comment表失败：%w", err)
	}
	return nil
}
