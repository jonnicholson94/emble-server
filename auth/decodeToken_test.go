package auth

import (
	"emble-server/utils"
	"os"
	"testing"
)

func TestDecodeToken(t *testing.T) {

	os.Setenv("DOTENV_PATH", "../.env")

	token, err := utils.CreateToken(1234, "Test", "Name")

	if err != nil {
		t.Errorf("Failed creating the token")
	}

	uid, err := DecodeTokenId(token)

	if err != nil {
		t.Errorf(err.Error())
	}

	expected := 1234

	if uid != float64(expected) {
		t.Error("User ID and expected do not match")
	}

}
