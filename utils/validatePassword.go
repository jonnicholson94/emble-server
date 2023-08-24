package utils

import "regexp"

func ValidatePassword(password string) bool {

	// Validate to check if password is correct length
	if len(password) < 6 || len(password) > 20 {
		return false
	}

	checkSpecialCharacters := regexp.MustCompile(`[^a-zA-Z0-9]`)
	if !checkSpecialCharacters.MatchString(password) {
		return false
	}

	checkForNumbers := regexp.MustCompile(`\d`)
	if !checkForNumbers.MatchString(password) {
		return false
	}

	return true

}
