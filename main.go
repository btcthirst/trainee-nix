package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/btcthirst/trainee-nix/docs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Comments struct use to database get/set data(front)
type Comments struct {
	PostID uint64 `json:"postId" gorm:"post_id"`
	ID     uint64 `json:"id" gorm:"id"`
	Name   string `json:"name" gorm:"name"`
	Email  string `json:"email" gorm:"email"`
	Body   string `json:"body" gorm:"body"`
}

// Posts struct use to database get/set data(front)
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
type Mess struct {
	Message interface{}
}

var (
	DB                 *gorm.DB
	PresentationFormat = "JSON" // "XML"
)

///////////////////service///////////////////////
func initSettings() string {
	// get env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Erorr load .env file")
	}
	setDB := SettingsDB{
		host:     os.Getenv("host"),
		port:     os.Getenv("port"),
		database: os.Getenv("database"),
		user:     os.Getenv("user"),
		password: os.Getenv("password"),
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

func ToHTML(c echo.Context, statusCode int, m Mess) error {

	if PresentationFormat == "XML" {
		res := c.XML(statusCode, m)
		return res
	} else {
		res := c.JSON(statusCode, m)
		return res
	}

}

///////////////////end service///////////////////////

//////////////////CRUD`s/////////////////////////////
/* Posts */
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

func updatePost(p Posts) error {
	var isP Posts
	if err := DB.Where("id=?", p.ID).First(&isP).Error; err != nil {
		return err
	}
	DB.Save(&p)
	return nil
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

func updateComment(c Comments) error {
	var isC Comments

	if err := DB.Where("id=?", c.ID).First(&isC).Error; err != nil {
		return err
	}
	DB.Save(&c)
	return nil
}

func deleteComment(id uint64) error {
	err := DB.Where("id=?", id).Delete(&Comments{}).Error
	if err != nil {
		return err
	}
	return nil
}

//////////////////end CRUD`s/////////////////////////

//////////////////// handlers /////////////////////////

// heloP godoc
// @Summary Show the hello message
// @Tags root
// @Produce json
// @Produce xml
// @Success 200 {object} Mess
// @Router / [get]
func helloP(c echo.Context) error {
	/* u := map[string]interface{}{
		"message": "hello page",
	} */
	m := Mess{
		Message: "Hello",
	}

	return ToHTML(c, http.StatusOK, m)

}

/// posts handlers

// postPosts godoc
// @Summary Create post
// @Tags posts
// @Param userId path int true "Posts.UserID"
// @Param title path string true "Posts.Title"
// @Param body path string true "Posts.Body"
// @Accept json
// @Produce json
// @Produce xml
// @Success 201 {object} Mess
// @Failure 400 {object} Mess
// @Router /posts/ [post]
func postPosts(c echo.Context) error {
	var (
		post Posts
		m    Mess
	)

	err := c.Bind(&post)
	if err != nil {
		m.Message = "empty post method"

		return ToHTML(c, http.StatusBadRequest, m)

	} else {
		err = createPost(post)
		if err != nil {
			m.Message = "post not created"

			return ToHTML(c, http.StatusBadRequest, m)

		} else {
			m.Message = "post created"

			return ToHTML(c, http.StatusCreated, m)
		}
	}

}

// getPostsBy godoc
// @Summary Get post by id
// @Tags posts
// @Param id path int true "Posts.ID"
// @Accept json
// @Produce json
// @Produce xml
// @Success 200 {object} Mess
// @Failure 400 {object} Mess
// @Router /posts/:id [get]
func getPostsBy(c echo.Context) error {
	var m Mess

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		m.Message = "wrong id"

		return ToHTML(c, http.StatusBadRequest, m)

	}
	post := getPost(id)

	m.Message = post

	return ToHTML(c, http.StatusOK, m)

}

// getPostsAll godoc
// @Summary Get all posts
// @Tags posts
// @Produce json
// @Produce xml
// @Success 200 {object} Mess
// @Router /posts/ [get]
func getPostsAll(c echo.Context) error {

	posts := getAllPosts()
	m := Mess{
		Message: posts,
	}

	return ToHTML(c, http.StatusOK, m)

}

// putPosts godoc
// @Summary update post
// @Tags posts
// @Param id path int true "Posts.ID"
// @Param userId path int true "Posts.UserID"
// @Param title path string true "Posts.Title"
// @Param body path string true "Posts.Body"
// @Accept json
// @Produce json
// @Produce xml
// @Success 200 {object} Mess
// @Failure 400 {object} Mess
// @Router /posts/:id [put]
func putPosts(c echo.Context) error {
	var (
		post Posts
		m    Mess
	)
	err := c.Bind(&post)
	if err != nil {
		m.Message = "empty put"

		return ToHTML(c, http.StatusBadRequest, m)

	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		m.Message = "wrong id"

		return ToHTML(c, http.StatusBadRequest, m)

	}
	post.ID = id

	err = updatePost(post)
	if err != nil {
		m.Message = fmt.Sprint(err)

		return ToHTML(c, http.StatusBadRequest, m)

	}
	m.Message = "post updated"

	return ToHTML(c, http.StatusOK, m)

}

// deletePosts godoc
// @Summary delete post
// @Tags posts
// @Param id path int true "Posts.ID"
// @Produce json
// @Produce xml
// @Success 200 {object} Mess
// @Failure 400 {object} Mess
// @Router /posts/:id [delete]
func deletePosts(c echo.Context) error {
	var (
		id  uint64
		err error
		m   Mess
	)
	id, err = strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		m.Message = "wrong id"

		return ToHTML(c, http.StatusBadRequest, m)

	}

	err = deletePost(id)
	if err != nil {
		m.Message = fmt.Sprint(err)

		return ToHTML(c, http.StatusBadRequest, m)

	}
	res := fmt.Sprintf("deleted post with id %d", id)
	m.Message = res

	return ToHTML(c, http.StatusOK, m)

}

/// end posts handlers

/// comments handlers

// postComments godoc
// @Summary Create comment
// @Tags comment
// @Param postId path int true "Comments.PostID"
// @Param name path string true "Comments.Name"
// @Param email path string true "Comments.Email"
// @Param body path string true "Comments.Body"
// @Accept json
// @Produce json
// @Produce xml
// @Success 201 {object} Mess
// @Failure 400 {object} Mess
// @Router /comments/ [post]
func postComments(c echo.Context) error {
	var (
		comment Comments
		m       Mess
	)
	err := c.Bind(&comment)
	if err != nil {
		m.Message = "empty comment"

		return ToHTML(c, http.StatusBadRequest, m)

	} else {
		err := createComment(comment)
		if err != nil {
			m.Message = "comment not created"

			return ToHTML(c, http.StatusBadRequest, m)

		}
		m.Message = "comment created"

		return ToHTML(c, http.StatusCreated, m)

	}
}

// getCommentsAll godoc
// @Summary Get all comments
// @Tags comment
// @Produce json/xml
// @Success 200 {object} Mess
// @Router /comments/ [get]
func getCommentsAll(c echo.Context) error {
	comments := getAllComments()
	m := Mess{
		Message: comments,
	}

	return ToHTML(c, http.StatusOK, m)

}

// getCommentsBy godoc
// @Summary Get comment by id
// @Tags comment
// @Param id path int true "Comments.ID"
// @Accept json
// @Produce json
// @Produce xml
// @Success 200 {object} Mess
// @Failure 400 {object} Mess
// @Router /comments/:id [get]
func getCommentsBy(c echo.Context) error {
	var (
		id  uint64
		err error
		m   Mess
	)
	id, err = strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		m.Message = "wrong id"

		return ToHTML(c, http.StatusBadRequest, m)

	}
	comment := getComment(id)
	m.Message = comment

	return ToHTML(c, http.StatusOK, m)

}

// putComments godoc
// @Summary Update comment
// @Tags comment
// @Param id path int true "Comments.ID"
// @Param postId path int true "Comments.PostID"
// @Param name path string true "Comments.Name"
// @Param email path string true "Comments.Email"
// @Param body path string true "Comments.Body"
// @Accept json
// @Produce json
// @Produce xml
// @Success 200 {object} Mess
// @Failure 400 {object} Mess
// @Router /comments/:id [put]
func putComments(c echo.Context) error {
	var (
		comment Comments
		m       Mess
		id      uint64
		err     error
	)
	err = c.Bind(&comment)
	if err != nil {
		m.Message = "empty put"

		return ToHTML(c, http.StatusBadRequest, m)

	}
	id, err = strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		m.Message = "wrong id"

		return ToHTML(c, http.StatusBadRequest, m)

	}
	comment.ID = id
	err = updateComment(comment)
	if err != nil {

		m.Message = fmt.Sprint(err)

		return ToHTML(c, http.StatusBadRequest, m)

	}
	m.Message = "comment updated"

	return ToHTML(c, http.StatusOK, m)

}

// deleteComments godoc
// @Summary delete comment
// @Tags comment
// @Param id path int true "Comments.ID"
// @Accept json
// @Produce json
// @Produce xml
// @Success 200 {object} Mess
// @Failure 400 {object} Mess
// @Router /comments/:id [delete]
func deleteComments(c echo.Context) error {
	var (
		id  uint64
		err error
		m   Mess
	)
	id, err = strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		m.Message = "wrong id"

		return ToHTML(c, http.StatusBadRequest, m)

	}

	err = deleteComment(id)
	if err != nil {
		m.Message = fmt.Sprint(err)

		return ToHTML(c, http.StatusBadRequest, m)

	}
	res := fmt.Sprintf("deleted comment with id %d", id)
	m.Message = res

	return ToHTML(c, http.StatusOK, m)
}

// end comments handlers

//////////////////// end handlers /////////////////////

// @title MY Example API
// @version 1.0
// @description This is my server for practice.

// @contact.name btcthirst
// @contact.url ...
// @contact.email btcthirst@gmail.com

// @host localhost:8181
// @BasePath /
func main() {
	DB = initDB()
	migrator()
	e := echo.New()

	e.GET("/", helloP)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	posts := e.Group("/posts")
	comments := e.Group("/comments")
	//posts routers
	posts.GET("/", getPostsAll)
	posts.GET("/:id", getPostsBy)
	posts.POST("/", postPosts)
	posts.PUT("/:id", putPosts)
	posts.DELETE("/:id", deletePosts)
	//comments routers
	comments.GET("/", getCommentsAll)
	comments.GET("/:id", getCommentsBy)
	comments.POST("/", postComments)
	comments.PUT("/:id", putComments)
	comments.DELETE("/:id", deleteComments)

	port := os.Getenv("PORT_SERVER")
	e.Logger.Fatal(e.Start(port))

}
