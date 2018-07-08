package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/superjinjo/catalyze-go/users"
)

type UserJSON struct {
	ID        int    `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Color     string `json:"color"`
}

//UsersController handles requests for /user
type UsersController struct {
	Users *users.Repository
}

//HandlePost creates a new user
func (con *UsersController) HandlePost(w http.ResponseWriter, r *http.Request) {
	var info UserJSON
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		WriteResponse(w, 400, Error{"Malformed JSON", "The JSON in the response body was malformed."})
		return
	}

	if info.Username == "" || info.Password == "" || info.Firstname == "" || info.Lastname == "" || info.Color == "" {
		WriteResponse(w, 400, Error{"Input Error", "Missing required fields."})
		return
	}

	user, err := con.Users.InsertUser(info.Username, info.Password, info.Firstname, info.Lastname, info.Color)
	logger.Println(info)
	if err != nil {
		WriteResponse(w, 500, Error{"Create Error", "An unknown error occured when creating the user."})
	} else {
		WriteResponse(w, 201, UserJSON{
			ID:        user.ID,
			Username:  user.Username,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Color:     user.Color,
		})
	}
}

//HandleGet gets the specified user from the route
func (con *UsersController) HandleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, err := con.Users.GetUser(vars["username"])
	logger.Println(user, err)
	if err != nil {
		WriteResponse(w, 500, Error{"Get Error", "An unknown error occured when getting the user."})
	} else {
		WriteResponse(w, 200, UserJSON{
			ID:        user.ID,
			Username:  user.Username,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Color:     user.Color,
		})
	}
}

//HandlePut updates the user specified in the route
func (con *UsersController) HandlePut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var info UserJSON

	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		WriteResponse(w, 400, Error{"Malformed JSON", "The JSON in the response body was malformed."})
		return
	}

	if info.Firstname == "" || info.Lastname == "" || info.Color == "" {
		WriteResponse(w, 400, Error{"Input Error", "Missing required fields."})
		return
	}

	user, err := con.Users.UpdateUser(vars["username"], info.Firstname, info.Lastname, info.Color)
	logger.Println(info)
	if err != nil {
		WriteResponse(w, 500, Error{"Update Error", "An unknown error occured when updating the user."})
	} else {
		WriteResponse(w, 200, UserJSON{
			ID:        user.ID,
			Username:  user.Username,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Color:     user.Color,
		})
	}
}

//HandleDelete removes a user from the database
func (con *UsersController) HandleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := con.Users.DeleteUser(vars["username"]); err != nil {
		WriteResponse(w, 500, Error{"Delete Error", "Problem deleting user."})
	} else {
		data := map[string]string{"message": "User was deleted successfully"}
		WriteResponse(w, 200, data)
	}
}
