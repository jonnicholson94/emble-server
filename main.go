package main

import (
	"emble-server/auth"
	"emble-server/utils"
	"net/http"
)

func main() {

	utils.Initialise()

	http.HandleFunc("/create-user", auth.CreateUser)
	http.HandleFunc("/sign-in", auth.SignIn)
	http.HandleFunc("/validate-user", utils.ValidateToken)

	http.ListenAndServe(":8080", nil)

	utils.GetDB().Close()
}
