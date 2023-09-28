package waitingList

import (
	"bytes"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/resendlabs/resend-go"
)

type CustomError struct {
	Message string `json:"message"`
	Status  int32  `json:"status"`
}

type SignUpDetails struct {
	Email string `json:"email"`
	Code  int32  `json:"code"`
}

func JoinBeta(w http.ResponseWriter, r *http.Request) {

	var signUp SignUpDetails
	currentTime := time.Now()

	decodeErr := json.NewDecoder(r.Body).Decode(&signUp)

	if decodeErr != nil {
		customErr := CustomError{
			Message: "Error decoding body",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	db := utils.GetDB()

	insertQuery := "INSERT INTO beta (email, code, timestamp) VALUES ($1, $2, $3)"

	_, err := db.Exec(insertQuery, signUp.Email, signUp.Code, currentTime)

	if err != nil {
		customErr := CustomError{
			Message: "There's been a problem signing you up. You might have already joined.",
			Status:  http.StatusBadRequest,
		}

		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	tmpl, err := template.ParseFiles("templates/betaWelcome.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// Render the template
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, nil)
	if err != nil {
		fmt.Println(err)
		customErr := CustomError{
			Message: "Failed to execute HTML template",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Emble <info@emble.app>",
		To:      []string{signUp.Email},
		Html:    tpl.String(),
		Subject: "Thanks for signing up!",
	}

	_, sendErr := client.Emails.Send(params)

	if sendErr != nil {
		fmt.Println(sendErr)
		customErr := CustomError{
			Message: "Failed to send email",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	res, _ := json.Marshal("Successfully signed up")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
