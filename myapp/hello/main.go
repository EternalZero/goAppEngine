package main
//package hello

import (
	"net/http"
	"html/template"
	"log"
	"github.com/nu7hatch/gouuid"
	//"strings"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type user struct{
	name string
	age string
}
var htmlTest *template.Template
var htmlTest2 *template.Template

//func init() {
func main(){
	var err error

	htmlTest, err = template.ParseFiles("login.html")

	if(err != nil){
		log.Panic(err)
	}

	htmlTest2, err = template.ParseFiles("postlogin.html")
	if(err != nil){
		log.Panic(err)
	}


	http.HandleFunc("/", UUIDCookie)

	http.HandleFunc("/login/", handler)
	http.HandleFunc("/login/postlogin.html", postLogin)
	http.ListenAndServe(":9090", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var err error

	err = htmlTest.Execute(w,nil)
	if(err != nil){
		log.Panic(err)
	}
}

func postLogin(w http.ResponseWriter, r * http.Request){
	var err error

	cookie, err := r.Cookie("session-fino")

	if(err == http.ErrNoCookie){
		UUIDCookie(w,r)
	}


	currUser := user{
		name: r.FormValue("name"),
		age: r.FormValue("age"),
	}


	userJSON, err := json.Marshal(currUser)
	if(err != nil){
		log.Panic(err)
	}

	//I DOnt think the marshal is working correctly all i get is base64 for {}
	fmt.Fprint(w, userJSON)

	userJSONString := base64.StdEncoding.EncodeToString(userJSON)


	//uuid := strings.Split(cookie.Value, "|")[0]
	cookie.Value = cookie.Value + "|" + userJSONString

	http.SetCookie(w, cookie)

	htmlTest2.Execute(w,nil)
	if(err != nil){
		log.Panic(err)
	}
}

func UUIDCookie(w http.ResponseWriter, r * http.Request){


	//checking if cookie UUID exists
	cookie, err := r.Cookie("session-fino")

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