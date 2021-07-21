package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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

func main() {
	p := new([]Posts)
	userID := 7
	myURL := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts?userId=%d", userID)
	resp, err := http.Get(myURL)
	checkout(err)

	body, err := ioutil.ReadAll(resp.Body)
	checkout(err)

	err = json.Unmarshal(body, &p)
	checkout(err)
	fmt.Println(p)
}
