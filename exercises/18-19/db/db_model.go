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
	UserID   uint `gorm:"not null"`
	User     User
	Comments []Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint `gorm:"not null"`
	Post    Post
	UserID  uint `gorm:"not null"`
	User    User
}

var DB *gorm.DB

func QueryAllUserPostAndComment(userId uint) []Post {
	var posts []Post
	DB.Model(&Post{}).Preload("Comments").Where("user_id = ?", userId).Find(&posts)
	return posts
}

func QueryPostCommentAdnCommentUser(postId uint) ([]Comment, map[uint]string) {
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
func QueryUserPostAndComment(userId uint) ([]Post, map[uint][]Comment) {
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
func CreateUser(user *User) uint {
	DB.Create(user)
	return user.ID
}
func CreatePost(post *Post) (uint, error) {
	userId := post.UserID
	var user []User
	DB.Find(&user, userId)
	if len(user) == 0 {
		return 0, fmt.Errorf("没有查到对应的用户，请注册：%d", userId)
	}
	DB.Create(post)
	return post.ID, nil
}

func CreateComment(comment *Comment) (uint, error) {
	postId := comment.PostID
	userID := comment.UserID
	var posts []Post
	var users []User
	DB.Find(&posts, postId)
	if len(posts) == 0 {
		return 0, fmt.Errorf("不能发表没有文章的评论")
	}
	DB.Find(&users, userID)
	if len(users) == 0 {
		return 0, fmt.Errorf("用户id不对，请注册")
	}
	DB.Create(&comment)
	return comment.ID, nil
}

func InitDB(dsn string) error {
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
