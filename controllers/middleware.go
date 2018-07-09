package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/superjinjo/catalyze-go/auth"
	"github.com/superjinjo/catalyze-go/users"
)

//Middleware provides functions to check authentication and authorization
type Middleware struct {
	Auth  *auth.Repository
	Users *users.Repository
}

//GetAuthUser returns the user id of the authorized user
func (authMW *Middleware) GetAuthUser(r *http.Request) (int, error) {
	var token string

	token = r.Header.Get("Authorization")
	if token != "" {
		if userID, err := authMW.Auth.CheckToken(token); err == nil {
			return userID, nil
		} else {
			return 0, Error{"Invalid Token", "The token you provided is invalid or expired."}
		}
	}

	return 0, Error{"No Token", "You did not provide an authentication token with your request."}
}

//IsAuthenticated checks if user has valid token in HTTP request Authorization header
func (authMW *Middleware) IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	_, err := authMW.GetAuthUser(r)

	if err != nil {
		WriteResponse(w, 401, err)
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
		userID, _ := authMW.GetAuthUser(r)
		authuser, err := authMW.Users.GetUsername(userID)
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
		WriteResponse(w, status, userErr)
	}
}
