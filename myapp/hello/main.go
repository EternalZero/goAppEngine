package main
//package hello

import (
	"net/http"
	"html/template"
	"log"
	"github.com/nu7hatch/gouuid"
)

var htmlTest *template.Template
var htmlTest2 *template.Template
//func init() {
func main(){
	var err error

	htmlTest, err = template.ParseFiles("test.html")

	if(err != nil){
		log.Panic(err)
	}

	htmlTest2, err = template.ParseFiles("test2.html")
	if(err != nil){
		log.Panic(err)
	}


	http.HandleFunc("/", UUIDCookie)

	http.HandleFunc("/feels/", handler)
	http.HandleFunc("/feels/test2.html", doIcare)
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

func UUIDCookie(w http.ResponseWriter, req * http.Request){


	//checking if cookie UUID exists
	cookie, err := req.Cookie("session-fino")

	//if the cookie does not exist we make a new and assign it a UUID using the imported gouuid thingy
	if(err == http.ErrNoCookie){
		//we ignore the error that the NewV4 function returns. we make a randomly generated UUID
		uuid, _ := uuid.NewV4()

		//we make a new cookie using a composite literal and give its address to a pointer
		cookie = &http.Cookie{
			Name: "session-fino",
			Value: uuid.String(),
			HttpOnly: true,
			//can't use secure cause we don't have https available
			// /Secure: true,
		}

		//we set the cookie on the users pc
		http.SetCookie(w, cookie)
	}
}