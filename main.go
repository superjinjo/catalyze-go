package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/superjinjo/catalyze-go/auth"
	"github.com/superjinjo/catalyze-go/users"
	"github.com/urfave/negroni"
)

func helloworld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello World"})
}

func setupAuthRoutes(n *negroni.Negroni, middleware *auth.Middleware, router *mux.Router, con *auth.Controller) {
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

func setupUserRoutes(n *negroni.Negroni, middleware *auth.Middleware, router *mux.Router, con *auth.Controller) {
	// nauth := n.With(
	// 	negroni.HandlerFunc(middleware.IsAuthenticated),
	// 	negroni.HandlerFunc(middleware.CanManageUser),
	// )

	//POST /user

	//GET /user/{username}

	//PUT /user/{username}

	//DELETE /user/{username}
}

func main() {

	userRepo := &users.Repository{DBuser: "catalyze", DBpw: "abcd1234", DBhost: "localhost", DBname: "catalyze"}
	defer userRepo.Close()

	authRepo := &auth.Repository{DBuser: "catalyze", DBpw: "abcd1234", DBhost: "localhost", DBname: "catalyze"}
	defer authRepo.Close()

	middleware := &auth.Middleware{Auth: authRepo, Users: userRepo}
	//userMW := negroni.HandlerFunc(middleware.CanManageUser)

	router := mux.NewRouter()

	n := negroni.Classic()

	router.Handle("/", n.With(
		negroni.Wrap(http.HandlerFunc(helloworld)),
	)).Methods("GET")

	authCon := &auth.Controller{Auth: authRepo, Users: userRepo}
	setupAuthRoutes(n, middleware, router, authCon)

	http.ListenAndServe(":3000", router)

	// data, err := repo.InsertUser("barney", "abcd1234", "Barney", "Rubles", "brown")
	// fmt.Printf("data: %v, err: %v\n", data, err)

	// if err == nil {
	// 	id, err := repo.CheckCredentials("barney", "abcd1234")
	// 	fmt.Printf("data: %v, err: %v\n", id, err)

	// 	if err == nil {
	// 		getdata, err := repo.GetUser("barney")
	// 		fmt.Printf("data: %v, err: %v\n", getdata, err)
	// 	}
	// }
}
