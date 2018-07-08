package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/superjinjo/catalyze-go/users"
)

//Error stores authentication error information and handles response
type Error struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func handleError(w http.ResponseWriter, status int, err Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

//Middleware provides functions to check authentication and authorization
type Middleware struct {
	Auth   Repository
	Users  users.Repository
	userID int
	token  string
}

//IsAuthenticated checks if user has valid token in HTTP request Authorization header
func (auth *Middleware) IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var token, title, message string
	var userID int

	token = r.Header.Get("Authorization")
	if token != "" && auth.token != token {
		if userID, err := auth.Auth.CheckToken(token); err == nil {
			auth.userID = userID
			auth.token = token
		} else {
			title = "Invalid Token"
			message = "The token you provided is invalid or expired."
		}
	} else if token == "" {
		title = "No Token"
		message = "You did not provide an authentication token with your request."
	}

	if userID == 0 {
		handleError(w, 401, Error{title, message})
	} else {
		next(w, r)
	}
}

//CanManageUser checks to see if the username in the URI matches that of the user.
func (auth *Middleware) CanManageUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	var userErr Error
	var status int

	if username, ok := vars["username"]; ok == true {
		authuser, err := auth.Users.GetUsername(auth.userID)
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
		handleError(w, status, userErr)
	}
}
