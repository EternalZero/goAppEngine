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
	"crypto/hmac"
	"crypto/sha256"
	"strings"
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
	http.HandleFunc("/checkCookie", verifyMessage)
	http.HandleFunc("/corruptCookie", corruptCookie)
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
		cookie.Value = cookie.Value + "|" + getHMAC(cookie.Value);
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


// I EDITED THIS CODE FROM my PREVIOUS HW
// MAKE SURE YOU KNOW WHAT I CHANGED!! its down there with the indexes being hmaced

//this function verifies the cookie by comparing the message recieved and checking that it matches the
//hmac code created by our hmac function
func verifyMessage(w http.ResponseWriter, req * http.Request){

	//asking for the specific cookie
	cookie, err := req.Cookie("session-fino")

	//if the cookie doesn't exist we offend the user
	if(err == http.ErrNoCookie){
		fmt.Fprint(w, "You didn't give me a cookie. jerk")
		return
	}

	//if the cookie exists we split the message and the hmac code
	//index 0 is the message and index 1 is the code
	//****HEY HEY HEY HEY!!!! So since this time we are working with a uuid | Data | hmac
	//we use index 0 and 1 when hmacing our message
	messageAndHMAC := strings.Split((cookie.Value), "|")

	//fmt.Fprintln(w,"message and hmac: ", messageAndHMAC)

	//we hash the message we got and save it
	verifier := getHMAC(messageAndHMAC[0] + messageAndHMAC[1])

	//fmt.Fprintln(w, "verifier:", verifier)
	//fmt.Fprintln(w, "hmac key from cookie:", messageAndHMAC[1])

	//we compare the message we calculated to the one we are supposed to match
	//**** REMEMBER WE HAVE MORE DATA. we are working with UUID | DATA | HMAC
	if(messageAndHMAC[2] != verifier){
		//if the message calculated hmac code doesn't match the one in the cookie
		//it means the user altered the cookie
		fmt.Fprint(w, "you messed up the cookie.... Way to go...")
		return
	}
	//if the message calculated hmac matches the one in the cookie then everything is fine
	fmt.Fprint(w, "Everything seems fine. Carry on..")


}

//returns a string that is hmac code created using a "secret" key and some data
func getHMAC(data string) string{
	//making a new hash using the sha256 algorithm and using the key provided
	h := hmac.New(sha256.New, []byte("123456"))

	//no need to print the hash value.... jk
	//fmt.Println(h, data)

	//format %x makes them all lowercase? or 2 characters per byte
	//sum returns a slice of bytes so it's converting the hash to a string
	//for us which we then return to whoever call this func
	return fmt.Sprintf("%x", h.Sum(nil))
}


func corruptCookie(w http.ResponseWriter, req * http.Request){
	//asking for the specific cookie
	cookie, err := req.Cookie("session-fino")

	//if the cookie doesn't exist we offend the user
	if(err == http.ErrNoCookie){
		fmt.Fprint(w, "You didn't give me a cookie. jerk")
		return
	}
	cookie.Value = cookie.Value + "corrupted"

	http.SetCookie(w,cookie)

	fmt.Fprint(w, "Your cookie has been corrupted :)")

}