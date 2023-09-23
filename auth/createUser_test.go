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

func setupCreateTest() func() {

	user := "test_email@go.com"

	utils.Initialise()

	reqBody, _ := json.Marshal(user)

	req, err := http.NewRequest("DELETE", "/delete-user", bytes.NewBuffer(reqBody))

	if err != nil {
		fmt.Println("Unable to delete user")
	}

	rr := httptest.NewRecorder()

	DeleteUser(rr, req)

	return func() {
		DeleteUser(rr, req)
	}

}

func TestCreateUser(t *testing.T) {

	teardown := setupCreateTest()

	defer teardown()

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
		t.Errorf(err.Error())
		fmt.Println(err)
	}

	rr := httptest.NewRecorder()

	CreateUser(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Code != 200 {
		t.Error(rr.Body.String())
	}

}
