package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"html/template"
	"log"
	"net/http"
)

type Soul struct {
	Id    int
	Title string `json:"title"`
	Hits  int    `json:"hits"`
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/soul?charset=utf8&parseTime=True&loc=Local&timeout=10ms")
	if err != nil {
		fmt.Printf("mysql connect error %v", err)
	}

	if db.Error != nil {
		fmt.Printf("database error %v", db.Error)
	}

	db.LogMode(true)
}

func main() {
	defer db.Close()

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/soul", soul)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe : ", err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	_ = t.Execute(w, nil)
}

func soul(w http.ResponseWriter, req *http.Request) {
	var soul Soul
	db.Model(&soul).Order("rand()").First(&soul)
	db.Model(&soul).UpdateColumn("hits", soul.Hits+1)
	soulJson, _ := json.Marshal(struct {
		Title string `json:"title"`
		Hits  int    `json:"hits"`
	}{
		Title: soul.Title,
		Hits:  soul.Hits,
	})
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(soulJson)
}
