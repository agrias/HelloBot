package api

import "testing"

func TestMakeDiscordAuthPOST(t *testing.T) {

	client_id := "462698059701420052"
	secret := "z976rWnj8rgxiKjOD_rj578GcCS3hrd5"

	request_json, err := GetJsonBytes(client_id, "OAuth2 Test", "", "", secret, nil)
	if (err != nil) {
		panic(err)
	}

	MakeDiscordAuthPOST(request_json)
}