//package main
package hello

import (
	"net/http"
	"html/template"
	"log"
)

var htmlTest *template.Template
var htmlTest2 *template.Template
func init() {
//func main(){
	var err error

	htmlTest, err = template.ParseFiles("test.html")

	if(err != nil){
		log.Panic(err)
	}

	htmlTest2, err = template.ParseFiles("test2.html")
	if(err != nil){
		log.Panic(err)
	}


	http.HandleFunc("/", handler)
	http.HandleFunc("/test2.html", doIcare)
	http.ListenAndServe(":9090", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var err error

	err = htmlTest.Execute(w,nil)
	if(err != nil){
		log.Panic(err)
	}
}

func doIcare(w http.ResponseWriter, r * http.Request){
	var err error

	htmlTest2.Execute(w,nil)
	if(err != nil){
		log.Panic(err)
	}
}