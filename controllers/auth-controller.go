package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/superjinjo/catalyze-go/auth"
	"github.com/superjinjo/catalyze-go/users"
)

//Tokenlife represents the number of seconds a token is good for (24 hours)
const Tokenlife int = 60 * 60 * 24

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//AuthController handles requests for /auth
type AuthController struct {
	Auth  *auth.Repository
	Users *users.Repository
}

//HandlePost generates a new authentication token
func (con *AuthController) HandlePost(w http.ResponseWriter, r *http.Request) {
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
func (con *AuthController) HandleDelete(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if err := con.Auth.DeleteToken(token); err != nil {
		writeResponse(w, 400, Error{"Delete Error", "Problem invalidating token. It might already be invalid."})
	} else {
		data := map[string]string{"message": "Token has been successfully invalidated"}
		writeResponse(w, 200, data)
	}
}
