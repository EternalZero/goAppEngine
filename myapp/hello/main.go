package hello

import (
	"net/http"
	"html/template"
	"log"
)

var htmlTest *template.Template

func init() {
	var err error

	htmlTest, err = template.ParseFiles("test.html", "test2.html")

	if(err != nil){
		log.Panic(err)
	}

	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var err error

	err = htmlTest.Execute(w,nil)
	if(err != nil){
		log.Panic(err)
	}


}