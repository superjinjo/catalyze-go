package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/superjinjo/catalyze-go/users"
)

var logger = log.New(os.Stdout, "http: ", log.LstdFlags)

//Tokenlife represents the number of seconds a token is good for (24 hours)
const Tokenlife int = 60 * 60 * 24

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Error stores authentication error information and handles response
type Error struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func writeResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

//Middleware provides functions to check authentication and authorization
type Middleware struct {
	Auth   *Repository
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

//Controller handles routing requests
type Controller struct {
	Auth  *Repository
	Users *users.Repository
}

//HandlePost generates a new authentication token
func (con *Controller) HandlePost(w http.ResponseWriter, r *http.Request) {
	var cred credentials
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		writeResponse(w, 400, Error{"Malformed JSON", "The JSON in the response body was malformed."})
		return
	}

	userID, err := con.Users.CheckCredentials(cred.Username, cred.Password)
	logger.Println(cred)
	logger.Println(con.Users.GetUser(cred.Username))
	if err != nil {
		writeResponse(w, 401, Error{"Invalid credentials", "The username and password were incorrect."})
		return
	}

	if token, err := con.Auth.CreateToken(userID, Tokenlife); err != nil {
		writeResponse(w, 500, Error{"Create Error", "An unknown error occured when creating the token."})
	} else {
		data := map[string]string{"token": token}
		writeResponse(w, 201, data)
	}
}

//HandleDelete invalidates the authentication in the Authorize header
func (con *Controller) HandleDelete(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if err := con.Auth.DeleteToken(token); err != nil {
		writeResponse(w, 400, Error{"Delete Error", "Problem invalidating token. It might already be invalid."})
	} else {
		data := map[string]string{"message": "Token has been successfully invalidated"}
		writeResponse(w, 200, data)
	}
}
