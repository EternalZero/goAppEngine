package main
//package hello

import (
	"net/http"
	"html/template"
	"log"
	"github.com/nu7hatch/gouuid"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type user struct{
	Name string
	Age string
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


	http.HandleFunc("/", getInfo)
	http.HandleFunc("/postlogin.html", postLogin)
	http.ListenAndServe(":9090", nil)
}


func postLogin(w http.ResponseWriter, r *http.Request){

	var err error

	//checking if cookie exists
	cookie, err := r.Cookie("session-fino")

	//if the cookie does not exist we make a new and assign it a UUID using the imported gouuid thingy
	if(err!=nil){
		fmt.Println("im making a cookie")
		//we ignore the error that the NewV4 function returns. we make a randomly generated UUID
		uuid, _ := uuid.NewV4()

		//we make a new cookie using a composite literal and give its address to a pointer
		cookie = &http.Cookie{
			Name: "session-fino",
			Value: uuid.String() + "|" + marshallUserType(r),
			HttpOnly: true,
			//can't use secure cause we don't have https available
			// /Secure: true,
		}

	fmt.Println(cookie)
	}
	http.SetCookie(w, cookie)
	//we set the cookie on the users pc

	err = htmlTest2.Execute(w,nil)
	if(err != nil){
		//log.Panic(err)
	}
}

func getInfo(w http.ResponseWriter, r * http.Request){

	var err error

	err = htmlTest.Execute(w,nil)
	if(err != nil){
		log.Panic(err)
	}
}

func marshallUserType(r * http.Request) string{

	currUser := user{
		Name: r.FormValue("name"),
		Age: r.FormValue("age"),
	}

	userJSON, err := json.Marshal(currUser)
	if(err != nil){
		log.Panic(err)
	}

	userJSONString := base64.StdEncoding.EncodeToString(userJSON)

	return userJSONString
}