//package main

package hello

import (
	"html/template"
	"net/http"
	//"log"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"strings"
)

type user struct {
	userName string
	password string
}

//we are adding user to our session data
//so now we can access the members of user through sessionData objects
type sessionData struct {
	thisUser  user
	loggedIn  bool
	loginFail bool
}

var tpl *template.Template

func init() {

	tpl = template.Must(template.ParseFiles("login.html", "postlogin.html"))

	http.HandleFunc("/", getInfo)
	http.HandleFunc("/postLogin", postLogin)
	http.HandleFunc("/logout", logout)

	http.HandleFunc("/checkUUID", verifyMessage)
	http.HandleFunc("/corruptCookie", corruptCookie)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func getInfo(w http.ResponseWriter, r *http.Request) {

	var err error
	context := appengine.NewContext(r)

	fmt.Fprintln(w, "printing stuff")
	err = tpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		log.Errorf(context, "error executing htmlTest template", err)
		http.Redirect(w, r, "/", 303)
	}
}

func postLogin(w http.ResponseWriter, r *http.Request) {

	var err error

	context := appengine.NewContext(r)

	if (!loginCheck(w, r)) {
		fmt.Println(w, "You got the password wrong :(")
		http.Redirect(w, r, "/", 303)
		return
	}

	//checking if cookie exists
	cookie, err := r.Cookie("session-uuid")

	//if the cookie does not exist we make a new and assign it a UUID using the imported gouuid thingy
	if err != nil {
		//fmt.Println("im making a cookie")
		//i ignore the error that the NewV4 function returns. we make a randomly generated UUID
		uuid, _ := uuid.NewV4()

		//we make a new cookie using a composite literal and give its address to a pointer
		cookie = &http.Cookie{
			Name:     "session-uuid",
			Value:    uuid.String(),
			HttpOnly: true,
			//can't use secure cause we don't have https available
			// /Secure: true,
		}

		fmt.Println(cookie)
		cookie.Value = cookie.Value + "|" + getHMAC(cookie.Value)
	}

	http.SetCookie(w, cookie)
	//we set the cookie on the users pc

	err = tpl.ExecuteTemplate(w, "postlogin.html", nil)
	if err != nil {
		log.Errorf(context, "error executing htmlTest2 template", err)
		http.Redirect(w, r, "/", 303)
		return
	}
}

func marshallUserType(w http.ResponseWriter, r *http.Request) string {

	context := appengine.NewContext(r)

	currUser := sessionData{
		thisUser: user{
			userName: r.FormValue("username"),
			password: r.FormValue("password"),
		},
		loggedIn:  false,
		loginFail: false,
	}

	userJSON, err := json.Marshal(currUser)
	if err != nil {
		log.Errorf(context, "error marshalling user object/instance", err)
		http.Redirect(w, r, "/", 303)
	}

	return string(userJSON)
}

func loginCheck(w http.ResponseWriter, r *http.Request) bool {

	cookie, err := r.Cookie("session-login")

	if err != http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:     "session-login",
			Value:    encodeToB64(marshallUserType(w, r)),
			HttpOnly: true,
			//Secure: true,
		}
	}

	jsonData, _ := base64.StdEncoding.DecodeString(cookie.Value)
	var currSession sessionData
	json.Unmarshal(jsonData, currSession)

	if r.Method == "POST" {
		password := r.FormValue("password")

		if password == "123456" {
			currSession.loggedIn = true
			currSession.loginFail = false
		} else {
			currSession.loggedIn = false
			currSession.loginFail = true
		}
	}

	loggedIn := currSession.loggedIn

	newJSON, _ := json.Marshal(currSession)

	cookie.Value = encodeToB64(string(newJSON))

	http.SetCookie(w, cookie)

	return loggedIn
}

func logout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session-login")

	if err != http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:     "session-login",
			Value:    encodeToB64(marshallUserType(w, r)),
			HttpOnly: true,
			//Secure: true,
		}
	} else {
		jsonData, _ := base64.StdEncoding.DecodeString(cookie.Value)
		var currSession sessionData
		json.Unmarshal(jsonData, currSession)
		currSession.loggedIn = false
	}

	http.SetCookie(w, cookie)
}

func encodeToB64(userJSON string) string {
	userJSONString := base64.StdEncoding.EncodeToString([]byte(userJSON))
	return userJSONString
}

// I EDITED THIS CODE FROM my PREVIOUS HW
// MAKE SURE YOU KNOW WHAT I CHANGED!! its down there with the indexes being hmaced

//this function verifies the cookie by comparing the message recieved and checking that it matches the
//hmac code created by our hmac function
func verifyMessage(w http.ResponseWriter, r *http.Request) {

	//asking for the specific cookie
	cookie, err := r.Cookie("session-uuid")

	//if the cookie doesn't exist we offend the user
	if err == http.ErrNoCookie {
		fmt.Fprint(w, "You didn't give me the uuid cookie. jerk1!")
		return
	}

	//if the cookie exists we split the message and the hmac code
	//index 0 is the message and index 1 is the code
	//****HEY HEY HEY HEY!!!! So since this time we are working with a uuid | Data | hmac
	//we use index 0 and 1 when hmacing our message
	messageAndHMAC := strings.Split((cookie.Value), "|")

	//fmt.Fprintln(w,"message and hmac: ", messageAndHMAC)

	//we hash the message we got and save it
	verifier := getHMAC(messageAndHMAC[0])

	//fmt.Fprintln(w, "verifier:", verifier)
	//fmt.Fprintln(w, "hmac key from cookie:", messageAndHMAC[1])

	//we compare the message we calculated to the one we are supposed to match
	//**** we are working with UUID | HMAC
	if messageAndHMAC[1] != verifier {
		//if the message calculated hmac code doesn't match the one in the cookie
		//it means the user altered the cookie
		fmt.Fprint(w, "you messed up the cookie.... Way to go...")
		return
	}
	//if the message calculated hmac matches the one in the cookie then everything is fine
	fmt.Fprint(w, "Everything seems fine. Carry on..")

}

//returns a string that is hmac code created using a "secret" key and some data
func getHMAC(data string) string {
	//making a new hash using the sha256 algorithm and using the key provided
	h := hmac.New(sha256.New, []byte("123456"))

	//no need to print the hash value.... jk
	//fmt.Println(h, data)

	//format %x makes them all lowercase? or 2 characters per byte
	//sum returns a slice of bytes so it's converting the hash to a string
	//for us which we then return to whoever call this func
	return fmt.Sprintf("%x", h.Sum(nil))
}

func corruptCookie(w http.ResponseWriter, r *http.Request) {
	//asking for the specific cookie
	cookie, err := r.Cookie("session-fino")

	//if the cookie doesn't exist we offend the user
	if err == http.ErrNoCookie {
		fmt.Fprint(w, "You didn't give me a cookie. jerk")
		return
	}
	cookie.Value = cookie.Value + "corrupted"

	http.SetCookie(w, cookie)

	fmt.Fprint(w, "Your cookie has been corrupted :)")

}
