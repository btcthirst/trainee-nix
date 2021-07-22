package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	DB *gorm.DB
)

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

func helloP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello page")
}

func main() {
	DB = initDB()
	migrator()

	http.HandleFunc("/", helloP)

	port := os.Getenv("PORT_SERVER")
	log.Fatal(http.ListenAndServe(port, nil))
}
