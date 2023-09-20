package auth

import (
	"emble-server/utils"
	"os"
	"testing"
)

func TestValidateToken(t *testing.T) {

	os.Setenv("DOTENV_PATH", "../.env")

	token, err := utils.CreateToken(1234, "Test", "Name")

	if err != nil {
		t.Errorf("Failed creating the token")
	}

	tokenErr := ValidateToken(token)

	if tokenErr != nil {
		t.Error(tokenErr.Error())
	}

}
