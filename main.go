package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/superjinjo/catalyze-go/auth"
	"github.com/superjinjo/catalyze-go/controllers"
	"github.com/superjinjo/catalyze-go/users"
	"github.com/urfave/negroni"
)

func helloworld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello World"})
}

func setupAuthRoutes(n *negroni.Negroni, middleware *controllers.Middleware, router *mux.Router, con *controllers.AuthController) {
	//POST /auth
	router.Handle("/auth", n.With(
		negroni.Wrap(http.HandlerFunc(con.HandlePost)),
	)).Methods("POST")

	//DELETE /auth
	router.Handle("/auth", n.With(
		negroni.HandlerFunc(middleware.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(con.HandleDelete)),
	)).Methods("DELETE")
}

func setupUserRoutes(n *negroni.Negroni, middleware *controllers.Middleware, router *mux.Router, con *controllers.UsersController) {
	nauth := n.With(
		negroni.HandlerFunc(middleware.IsAuthenticated),
		negroni.HandlerFunc(middleware.CanManageUser),
	)

	//POST /user
	router.Handle("/users", n.With(
		negroni.Wrap(http.HandlerFunc(con.HandlePost)),
	)).Methods("POST")

	//GET /user/{username}
	router.Handle("/users/{username}", nauth.With(
		negroni.Wrap(http.HandlerFunc(con.HandleGet)),
	)).Methods("GET")

	//PUT /user/{username}
	router.Handle("/users/{username}", nauth.With(
		negroni.Wrap(http.HandlerFunc(con.HandlePut)),
	)).Methods("PUT")

	//DELETE /user/{username}
	router.Handle("/users/{username}", nauth.With(
		negroni.Wrap(http.HandlerFunc(con.HandleDelete)),
	)).Methods("DELETE")
}

func main() {

	userRepo := &users.Repository{DBuser: "catalyze", DBpw: "abcd1234", DBhost: "localhost", DBname: "catalyze"}
	defer userRepo.Close()

	authRepo := &auth.Repository{DBuser: "catalyze", DBpw: "abcd1234", DBhost: "localhost", DBname: "catalyze"}
	defer authRepo.Close()

	middleware := &controllers.Middleware{Auth: authRepo, Users: userRepo}

	router := mux.NewRouter()

	n := negroni.Classic()

	//GET /
	router.Handle("/", n.With(
		negroni.Wrap(http.HandlerFunc(helloworld)),
	)).Methods("GET")

	//Routes for /auth
	authCon := &controllers.AuthController{Auth: authRepo, Users: userRepo}
	setupAuthRoutes(n, middleware, router, authCon)

	//Routes for /users
	usersCon := &controllers.UsersController{Users: userRepo}
	setupUserRoutes(n, middleware, router, usersCon)

	http.ListenAndServe(":3000", router)
}
