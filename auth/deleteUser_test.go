package auth

import (
	"bytes"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupDeleteTest() {

	newUser := User{
		FirstName: "Test",
		LastName:  "Name",
		Email:     "test_email@go.com",
		Password:  "testpassword123!",
	}

	utils.Initialise()

	reqBody, _ := json.Marshal(newUser)

	req, err := http.NewRequest("POST", "/create-user", bytes.NewBuffer(reqBody))

	if err != nil {
		fmt.Println("Unable to create user")
	}

	rr := httptest.NewRecorder()

	CreateUser(rr, req)

}

func TestDeleteUser(t *testing.T) {

	setupDeleteTest()

	user := UserToDelete{
		Email: "test_email@go.com",
	}

	utils.Initialise()

	reqBody, _ := json.Marshal(user)

	req, err := http.NewRequest("DELETE", "/delete-user", bytes.NewBuffer(reqBody))

	if err != nil {
		t.Error(err.Error())
	}

	rr := httptest.NewRecorder()

	DeleteUser(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Code != 200 {
		t.Error(rr.Body.String())
	}

}
