package auth

import (
	"bytes"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/google/uuid"
	"github.com/resendlabs/resend-go"
)

type EnteredEmail struct {
	Email string `json:"email"`
}

func PasswordReset(w http.ResponseWriter, r *http.Request) {

	// Parse email from req body

	var email EnteredEmail

	err := json.NewDecoder(r.Body).Decode(&email)

	if err != nil {
		customErr := CustomError{
			Message: "Failed to decode body",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	// Generate reset token

	id := uuid.New()

	// Save token into database

	db := utils.GetDB()

	query := "UPDATE users SET reset_token = $1 WHERE email = $2;"

	_, insertErr := db.Exec(query, id, email.Email)

	if insertErr != nil {
		fmt.Println(insertErr.Error())
		customErr := CustomError{
			Message: "Error inserting to database",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	// Send email containing appropriate link

	os.Setenv("DOTENV_PATH", "../.env")

	var customURL string

	env := os.Getenv("ENVIRONMENT")

	if env == "development" {
		customURL = "http://localhost:3000/auth/reset-password?id=" + id.String()
	} else {
		customURL = "https://emble.app/auth/reset-password?id=" + id.String()
	}

	tmpl, err := template.ParseFiles("templates/auth/forgotPassword.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// Prepare data for the template
	data := struct {
		CustomURL string
	}{
		CustomURL: customURL,
	}

	// Render the template
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
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
		To:      []string{email.Email},
		Html:    tpl.String(),
		Subject: "Your password reset on emble",
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

	json, _ := json.Marshal("Successfully sent password reset email")

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
