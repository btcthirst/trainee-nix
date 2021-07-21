package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func checkout(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getComments(postID uint64) {
	comments := new([]Comments)
	myURL := fmt.Sprintf("https://jsonplaceholder.typicode.com/comments?postId=%d", postID)
	resp, err := http.Get(myURL)
	checkout(err)

	body, err := ioutil.ReadAll(resp.Body)
	checkout(err)

	err = json.Unmarshal(body, &comments)
	checkout(err)
	fmt.Println(comments)
}

func main() {
	posts := new([]Posts)
	userID := 7
	myURL := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts?userId=%d", userID)
	resp, err := http.Get(myURL)
	checkout(err)

	body, err := ioutil.ReadAll(resp.Body)
	checkout(err)

	err = json.Unmarshal(body, &posts)
	checkout(err)

	for _, p := range *posts {
		go getComments(p.ID)

	}
	var stop string
	_, err = fmt.Scan(&stop)
}
