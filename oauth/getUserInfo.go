package oauth

import "net/http"

func GetUserInfo(access_token string) (*http.Response, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+access_token)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil

}
