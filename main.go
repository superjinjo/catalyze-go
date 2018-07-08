package main

import (
	"fmt"

	"github.com/superjinjo/catalyze-go/users"
)

func main() {
	fmt.Println("hello world")

	repo := users.Repository{DBuser: "catalyze", DBpw: "abcd1234", DBhost: "localhost", DBname: "catalyze"}
	defer repo.Close()

	data, err := repo.InsertUser("barney", "abcd1234", "Barney", "Rubles", "brown")
	fmt.Printf("data: %v, err: %v\n", data, err)

	if err == nil {
		id, err := repo.GetID("barney", "abcd1234")
		fmt.Printf("data: %v, err: %v\n", id, err)

		if err == nil {
			getdata, err := repo.GetUser("barney")
			fmt.Printf("data: %v, err: %v\n", getdata, err)
		}
	}
}
