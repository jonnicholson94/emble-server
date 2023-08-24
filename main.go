package main

import (
	"emble-server/auth"
	"emble-server/crud"
	"emble-server/utils"
	"emble-server/ws"
	"net/http"
)

func main() {

	utils.Initialise()

	http.HandleFunc("/create-user", auth.CreateUser)
	http.HandleFunc("/sign-in", auth.SignIn)
	http.HandleFunc("/validate-user", utils.ValidateToken)
	http.HandleFunc("/create-research", crud.CreateResearch)

	http.HandleFunc("/ws", ws.Websocket)

	http.ListenAndServe(":8080", nil)

	utils.GetDB().Close()
}
