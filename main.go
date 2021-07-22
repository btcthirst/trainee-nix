package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

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
	DB                 *gorm.DB
	PresentationFormat = "JSON" //"XML"
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

func toHtml(w http.ResponseWriter, data interface{}, statusCode int) {
	switch PresentationFormat {
	case "JSON":
		w.Header().Set("Content-type", "application/json; charset=UTF8")
		w.WriteHeader(statusCode)
		err := json.NewEncoder(w).Encode(data)
		checkout(err)
	case "XML":
		w.Header().Set("Content-type", "application/xml; charset=UTF8")
		w.WriteHeader(statusCode)
		err := xml.NewEncoder(w).Encode(data)
		checkout(err)
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
		p.ID = 100500
		p.Title = "Phantom Post"
		p.Body = "you can see this message because there are no posts in the database"
		p.UserID = 100500
	}
	return p
}

func getAllPosts() []Posts {
	var posts []Posts
	r := DB.Find(&posts)
	if r.RowsAffected == 0 {
		p := Posts{
			ID:     100500,
			Title:  "Phantom Post",
			Body:   "you can see this message because there are no posts in the database",
			UserID: 100500,
		}
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
		c.ID = 100500
		c.Name = "Phantom Post"
		c.Email = "Phantom Email"
		c.Body = "you can see this message because there are no comments in the database"
		c.PostID = 100500
	}
	return c
}

func getAllComments() []Comments {
	var comments []Comments
	r := DB.Find(&comments)
	if r.RowsAffected == 0 {
		c := Comments{
			ID:     100500,
			Name:   "Phantom Post",
			Email:  "Phantom Email",
			Body:   "you can see this message because there are no comments in the database",
			PostID: 100500,
		}
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
func helloP(w http.ResponseWriter, r *http.Request) {
	toHtml(w, "hello page", http.StatusOK)
}

func postsPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postPosts(w, r)
	case http.MethodGet:
		getPosts(w, r)
	case http.MethodPut:
		putPosts(w, r)
	case http.MethodDelete:
		deletePosts(w, r)
	default:
		data := fmt.Sprintf("%v -there is no such method on this page", r.Method)
		toHtml(w, data, http.StatusMethodNotAllowed)
	}
}

///methods for posts page

func postPosts(w http.ResponseWriter, r *http.Request) {
	var post Posts
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		toHtml(w, err, http.StatusNoContent)
	} else {
		err = createPost(post)
		if err != nil {
			toHtml(w, err, http.StatusBadRequest)
		} else {
			toHtml(w, post, http.StatusCreated)
		}
	}

}

func getPosts(w http.ResponseWriter, r *http.Request) {
	slice := strings.Split(r.URL.String(), "/")
	lastEl := slice[len(slice)-1]
	if lastEl != "" {
		id, err := strconv.ParseUint(lastEl, 10, 64)
		if err != nil {
			toHtml(w, "wrong id", http.StatusBadRequest)
		}
		posts := getPost(id)
		toHtml(w, posts, http.StatusOK)
	} else {
		posts := getAllPosts()
		toHtml(w, posts, http.StatusOK)
	}
}

func putPosts(w http.ResponseWriter, r *http.Request) {
	var post Posts
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		toHtml(w, err, http.StatusNoContent)
	} else {
		slice := strings.Split(r.URL.String(), "/")
		lastEl := slice[len(slice)-1]
		if lastEl != "" {
			id, err := strconv.ParseUint(lastEl, 10, 64)
			if err != nil {
				toHtml(w, "wrong id", http.StatusBadRequest)
			}
			post.ID = id
			updatePost(post)
			toHtml(w, post, http.StatusOK)
		} else {
			toHtml(w, "no id", http.StatusBadRequest)
		}
	}

}

func deletePosts(w http.ResponseWriter, r *http.Request) {
	slice := strings.Split(r.URL.String(), "/")
	lastEl := slice[len(slice)-1]
	if lastEl != "" {
		id, err := strconv.ParseUint(lastEl, 10, 64)
		if err != nil {
			toHtml(w, "wrong id", http.StatusBadRequest)
		}
		err = deletePost(id)
		if err != nil {
			toHtml(w, err, http.StatusBadRequest)
		}
		toHtml(w, id, http.StatusOK)
	} else {
		toHtml(w, "no id", http.StatusBadRequest)
	}
}

///end methods for posts page

func commentsPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postComments(w, r)
	case http.MethodGet:
		getComments(w, r)
	case http.MethodPut:
		putComments(w, r)
	case http.MethodDelete:
		deleteComments(w, r)
	default:
		data := fmt.Sprintf("%v -there is no such method on this page", r.Method)
		toHtml(w, data, http.StatusMethodNotAllowed)
	}
}

///methods for comments page

func postComments(w http.ResponseWriter, r *http.Request) {
	var comment Comments
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		toHtml(w, err, http.StatusNoContent)
	} else {
		err = createComment(comment)
		if err != nil {
			toHtml(w, err, http.StatusBadRequest)
		} else {
			toHtml(w, comment, http.StatusCreated)
		}
	}

}

func getComments(w http.ResponseWriter, r *http.Request) {
	slice := strings.Split(r.URL.String(), "/")
	lastEl := slice[len(slice)-1]
	if lastEl != "" {
		id, err := strconv.ParseUint(lastEl, 10, 64)
		if err != nil {
			toHtml(w, "wrong id", http.StatusBadRequest)
		}
		comment := getComment(id)
		toHtml(w, comment, http.StatusOK)
	} else {
		comments := getAllComments()
		toHtml(w, comments, http.StatusOK)
	}
}

func putComments(w http.ResponseWriter, r *http.Request) {
	var comment Comments
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		toHtml(w, err, http.StatusNoContent)
	} else {
		slice := strings.Split(r.URL.String(), "/")
		lastEl := slice[len(slice)-1]
		if lastEl != "" {
			id, err := strconv.ParseUint(lastEl, 10, 64)
			if err != nil {
				toHtml(w, "wrong id", http.StatusBadRequest)
			}
			comment.ID = id
			updateComment(comment)
			toHtml(w, comment, http.StatusOK)
		} else {
			toHtml(w, "no id", http.StatusBadRequest)
		}
	}

}

func deleteComments(w http.ResponseWriter, r *http.Request) {
	slice := strings.Split(r.URL.String(), "/")
	lastEl := slice[len(slice)-1]
	if lastEl != "" {
		id, err := strconv.ParseUint(lastEl, 10, 64)
		if err != nil {
			toHtml(w, "wrong id", http.StatusBadRequest)
		}
		err = deleteComment(id)
		if err != nil {
			toHtml(w, err, http.StatusBadRequest)
		}
		toHtml(w, id, http.StatusOK)
	} else {
		toHtml(w, "no id", http.StatusBadRequest)
	}
}

///end methods for posts page

////////////////////end handlers/////////////////////
func main() {
	DB = initDB()
	migrator()

	http.HandleFunc("/", helloP)
	http.HandleFunc("/posts/", postsPage)
	http.HandleFunc("/comments/", commentsPage)

	port := os.Getenv("PORT_SERVER")
	log.Fatal(http.ListenAndServe(port, nil))
}
