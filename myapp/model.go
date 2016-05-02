package main

import ()

type User struct {
	Email    string
	UserName string `datastore:"-"`
	Password string `json:"-"`
}

type SessionData struct {
	User
	LoggedIn      bool
	LoginFail     bool
	//not sure if we need this viewinguser string
	//ViewingUser   string

}

