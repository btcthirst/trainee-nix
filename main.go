package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Comments struct {
	PostID uint64 `json:"postId" gorm:"post_id"`
	ID     uint64 `json:"id" gorm:"id"`
	Name   string `json:"name" gorm:"name"`
	Email  string `json:"email" gorm:"email"`
	Body   string `json:"body" gorm:"body"`
}

type Posts struct {
	UserID uint64 `json:"userId" gorm:"user_id"`
	ID     uint64 `json:"id" gorm:"id"`
	Title  string `json:"title" gorm:"title"`
	Body   string `json:"body" gorm:"body"`
}

type SettingsDB struct {
	host     string
	port     string
	database string
	user     string
	password string
}

var (
	DB *gorm.DB
)

///////////////////service///////////////////////
func initSettings() string {
	// get env variables
	err := godotenv.Load(".env")

	setDB := SettingsDB{
		host:     os.Getenv("host"),
		port:     os.Getenv("port"),
		database: os.Getenv("database"),
		user:     os.Getenv("user"),
		password: os.Getenv("password"),
	}

	if err != nil {
		log.Fatal("Erorr load .env file")
	}
	// use env vars
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", setDB.user, setDB.password, setDB.host, setDB.port, setDB.database)
	return dsn
}

func initDB() *gorm.DB {

	dsn := initSettings()
	//connect to mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	checkout(err)
	return db
}

func migrator() {
	if DB.Migrator().HasTable(&Posts{}) {
		if DB.Migrator().HasTable(&Comments{}) {
			DB.Migrator().DropTable(&Posts{}, &Comments{})
		}
		DB.Migrator().DropTable(&Posts{})
	}
	DB.Migrator().CreateTable(&Posts{}, &Comments{})
}

func checkout(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

///////////////////end service///////////////////////

//////////////////CRUD`s/////////////////////////////
///Posts
func createPost(p Posts) error {
	if err := DB.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

func getPost(id uint64) Posts {
	var p Posts
	DB.Where("id=?", id).First(&p)
	if p.ID == 0 {
		p.Title = "No such post"
		p.Body = "if you want to read this post maybe create it"
	}

	return p
}

func getAllPosts() []Posts {
	var posts []Posts
	DB.Find(&posts)
	if len(posts) == 0 {
		var p Posts
		p.Title = "No any posts"
		p.Body = "if you want to read some of these posts maybe post them"
		posts = append(posts, p)
	}

	return posts
}

func updatePost(p Posts) {
	DB.Save(&p)
}

func deletePost(id uint64) error {
	err := DB.Where("id=?", id).Delete(&Posts{}).Error
	if err != nil {
		return err
	}
	return nil
}

///Comments
func createComment(c Comments) error {

	if err := DB.Create(&c).Error; err != nil {
		return err
	}
	return nil
}

func getComment(id uint64) Comments {
	var c Comments
	DB.Where("id=?", id).First(&c)
	if c.ID == 0 {
		c.Name = "No such comment"
		c.Body = "create a comment and try again"
	}
	return c
}

func getAllComments() []Comments {
	var comments []Comments
	DB.Find(&comments)
	if len(comments) == 0 {
		var c Comments
		c.Name = "No any comments"
		c.Body = "create a comment and try again"
		comments = append(comments, c)
	}
	return comments
}

func updateComment(c Comments) {
	DB.Save(&c)
}

func deleteComment(id uint64) error {
	err := DB.Where("id=?", id).Delete(&Comments{}).Error
	if err != nil {
		return err
	}
	return nil
}

//////////////////end CRUD`s/////////////////////////

////////////////////handlers/////////////////////////
func helloP(c echo.Context) error {

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "hello page",
	})

}

///methods for posts page

func postPosts(c echo.Context) error {
	var post Posts
	err := c.Bind(&post)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "empty post method",
		})
	} else {
		err = createPost(post)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "post not created",
			})

		} else {
			return c.JSON(http.StatusCreated, map[string]interface{}{
				"message": "post created",
			})

		}
	}

}

func getPostsBy(c echo.Context) error {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "wrong id",
		})

	}
	post := getPost(id)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"post": post,
	})

}

func getPostsAll(c echo.Context) error {

	posts := getAllPosts()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"post": posts,
	})

}

func putPosts(c echo.Context) error {
	var post Posts
	err := c.Bind(&post)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "empty put",
		})
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "wrong id",
		})

	}
	post.ID = id

	updatePost(post)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "post updated",
	})

}

func deletePosts(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "wrong id",
		})

	}

	err = deletePost(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "wrong id",
		})
	}
	res := fmt.Sprintf("deleted post with id %d", id)
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"message": res,
	})
}

///end methods for posts page

///methods for comments page

func postComments(c echo.Context) error {
	var comment Comments
	err := c.Bind(&comment)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "empty comment",
		})
	} else {
		err := createComment(comment)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "comment not created",
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "comment created",
		})
	}
}

func getCommentsAll(c echo.Context) error {
	comments := getAllComments()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"comments": comments,
	})
}

func getCommentsBy(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "wrong id",
		})

	}
	comment := getComment(id)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"comment": comment,
	})

}

func putComments(c echo.Context) error {
	var comment Comments
	err := c.Bind(&comment)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "empty put",
		})
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "wrong id",
		})

	}
	comment.ID = id
	updateComment(comment)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "comment updated",
	})
}

func deleteComments(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "wrong id",
		})

	}

	err = deleteComment(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "wrong id",
		})
	}
	res := fmt.Sprintf("deleted comment with id %d", id)
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"message": res,
	})
}

//end methods for posts page

////////////////////end handlers/////////////////////
func main() {
	DB = initDB()
	migrator()
	e := echo.New()

	e.GET("/", helloP)

	posts := e.Group("/posts")
	comments := e.Group("/comments")
	//posts routers
	posts.GET("/", getPostsAll)
	posts.GET("/:id", getPostsBy)
	posts.POST("/", postPosts)
	posts.PUT("/:id", putPosts)
	posts.DELETE("/:id", deletePosts)

	comments.GET("/", getCommentsAll)
	comments.GET("/:id", getCommentsBy)
	comments.POST("/", postComments)
	comments.PUT("/:id", putComments)
	comments.DELETE("/:id", deleteComments)

	port := os.Getenv("PORT_SERVER")
	e.Logger.Fatal(e.Start(port))

}
