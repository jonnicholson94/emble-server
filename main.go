package main

import (
	"emble-server/auth"
	"emble-server/crud"
	"emble-server/utils"
	"emble-server/ws"
	"net/http"

	"github.com/rs/cors"
)

func main() {

	utils.Initialise()

	mux := http.NewServeMux()

	mux.HandleFunc("/create-user", auth.CreateUser)
	mux.HandleFunc("/sign-in", auth.SignIn)
	mux.HandleFunc("/create-research", crud.CreateResearch)
	mux.HandleFunc("/research", crud.FetchResearch)
	mux.HandleFunc("/single-research", crud.FetchSingleResearch)
	mux.HandleFunc("/create-question", crud.CreateQuestion)
	mux.HandleFunc("/edit-item", crud.EditItem)
	mux.HandleFunc("/update-question-order", crud.UpdateQuestionOrder)
	mux.HandleFunc("/add-options", crud.AddOption)
	mux.HandleFunc("/edit-question", crud.UpdateQuestion)
	mux.HandleFunc("/delete-research", crud.DeleteResearch)
	mux.HandleFunc("/create-comment", crud.CreateComment)
	mux.HandleFunc("/edit-comment", crud.EditComment)

	mux.HandleFunc("/ws", ws.Websocket)

	// Very important!!!
	// Essential to review the cors AllowAll before app is production ready.

	handler := cors.AllowAll().Handler(mux)

	http.ListenAndServe(":8080", handler)

	utils.GetDB().Close()
}
