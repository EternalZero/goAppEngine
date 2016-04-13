package myapp

import (
	//"fmt"
	//"net/http"
	//"google.golang.org/appengine"
	//"google.golang.org/appengine/log"
	//"io"
	//"encoding/base64"
	//"encoding/json"
	//"crypto/sha1"
	//"mime/multipart"
	//"strings"
	//"golang.org/x/net/context"
	//"google.golang.org/cloud/storage"
	"html/template"
	"net/http"
)

var template1 *template.Template

func init() {

	var err error

	template1, err = template.ParseFiles("voteForMe.html")
	if(err != nil){
		//log.Fatal("parsefiles: ", err)
		//use googles log.Errorf(faodsi)
	}

	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {

	var err error

	if(r.URL.Path != "/"){
		http.NotFound(w,r)
		return
	}

	template1, err = template.ParseFiles("voteForMe.html")
	if(err != nil){
		//log.Fatal("parsefiles: ", err)
		//use googles log.Errorf(faodsi)
	}

	err = template1.Execute(w, nil)

	if(err != nil){
		//log.Fatal("exexute: ", err)
		//user googles log.Errorf(fjka;f)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

}
