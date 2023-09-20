package questions

import (
	"bytes"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestCreateQuestion(t *testing.T) {

	researchUUID := uuid.New()
	questionUUID := uuid.New()
	optionUUID := uuid.New()

	newOptions := Option{
		OptionId:         optionUUID.String(),
		OptionContent:    "Test content",
		OptionQuestionID: questionUUID.String(),
		OptionIndex:      1,
		OptionResearchID: researchUUID.String(),
	}

	newQuestion := NewQuestion{
		QuestionId:         questionUUID.String(),
		QuestionTitle:      "Question title",
		QuestionType:       "Single select",
		QuestionOptions:    []Option{newOptions},
		QuestionResearchId: researchUUID.String(),
		QuestionIndex:      1,
	}

	os.Setenv("DOTENV_PATH", "../.env")

	utils.Initialise()

	reqBody, _ := json.Marshal(newQuestion)

	req, err := http.NewRequest("POST", "/create-question", bytes.NewBuffer(reqBody))

	fmt.Println(req)

	if err != nil {
		t.Errorf(err.Error())
		fmt.Println(err)
	}

	token, err := utils.CreateToken(123, "Test", "Account")

	if err != nil {
		fmt.Println(err)
		t.Errorf(err.Error())
	}

	req.Header.Set("Authorization", token)

	rr := httptest.NewRecorder()

	fmt.Println(rr)

	CreateQuestion(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Code != 200 {
		t.Error("Status code doesn't match")
	}

}
