package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/superjinjo/catalyze-go/auth"
	"github.com/superjinjo/catalyze-go/users"
)

//Middleware provides functions to check authentication and authorization
type Middleware struct {
	Auth   *auth.Repository
	Users  *users.Repository
	userID int
	token  string
}

//IsAuthenticated checks if user has valid token in HTTP request Authorization header
func (authMW *Middleware) IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var token, title, message string

	token = r.Header.Get("Authorization")
	if token != "" && authMW.token != token {
		if userID, err := authMW.Auth.CheckToken(token); err == nil {
			authMW.userID = userID
			authMW.token = token
		} else {
			title = "Invalid Token"
			message = "The token you provided is invalid or expired."
		}
	} else if token == "" {
		title = "No Token"
		message = "You did not provide an authentication token with your request."
	}

	logger.Println(authMW)

	if authMW.userID == 0 {
		writeResponse(w, 401, Error{title, message})
	} else {
		next(w, r)
	}
}

//CanManageUser checks to see if the username in the URI matches that of the user.
func (authMW *Middleware) CanManageUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	var userErr Error
	var status int

	if username, ok := vars["username"]; ok == true {
		authuser, err := authMW.Users.GetUsername(authMW.userID)
		if err == nil && authuser == username {
			next(w, r)
		} else {
			userErr = Error{"Unauthorized", "You are not authorized to view or modify this user."}
			status = 403
		}
	} else {
		userErr = Error{"Bad request", "Request does not have username in URI."}
		status = 400
	}

	if status != 0 {
		writeResponse(w, status, userErr)
	}
}
