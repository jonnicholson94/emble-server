package main

import (
	"emble-server/auth"
	"emble-server/comments"
	"emble-server/options"
	"emble-server/questions"
	"emble-server/research"
	"emble-server/survey"
	"emble-server/utils"
	waitingList "emble-server/waiting-list"
	"emble-server/ws"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	os.Setenv("DOTENV_PATH", "./.env")

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	utils.Initialise()

	mux := http.NewServeMux()

	// Auth

	mux.HandleFunc("/create-user", auth.CreateUser)
	mux.HandleFunc("/delete-user", auth.DeleteUser)
	mux.HandleFunc("/sign-in", auth.SignIn)

	// Comments

	mux.HandleFunc("/create-comment", comments.CreateComment)
	mux.HandleFunc("/delete-comment", comments.DeleteComment)
	mux.HandleFunc("/edit-comment", comments.EditComment)

	// Options

	mux.HandleFunc("/create-option", options.CreateOption)
	mux.HandleFunc("/delete-option", options.DeleteOption)
	mux.HandleFunc("/edit-option", options.EditOption)

	// Questions

	mux.HandleFunc("/create-question", questions.CreateQuestion)
	mux.HandleFunc("/delete-question", questions.DeleteQuestion)
	mux.HandleFunc("/edit-question", questions.UpdateQuestion)
	mux.HandleFunc("/update-question-order", questions.UpdateQuestionOrder)

	// Research

	mux.HandleFunc("/create-research", research.CreateResearch)
	mux.HandleFunc("/delete-research", research.DeleteResearch)
	mux.HandleFunc("/edit-item", research.EditItem)
	mux.HandleFunc("/research", research.FetchResearch)
	mux.HandleFunc("/single-research", research.FetchSingleResearch)

	// Survey

	mux.HandleFunc("/survey", survey.FetchSurveyDetails)
	mux.HandleFunc("/create-response", survey.CreateResponse)

	// Waiting list

	mux.HandleFunc("/join-beta", waitingList.JoinBeta)

	// Web socket

	mux.HandleFunc("/ws", ws.Websocket)

	// Very important!!!
	// Essential to review the cors AllowAll before app is production ready.

	handler := cors.AllowAll().Handler(mux)

	port := os.Getenv("PORT")

	http.ListenAndServe(port, handler)

	utils.GetDB().Close()
}
