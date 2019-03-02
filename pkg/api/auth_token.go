package api

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
)



type DiscordAuthRequest struct {
	id 				string 		`json:"id"`
	name 			string 		`json:"name"`
	description	 	string 		`json:"description"`
	icon 			string 		`json:"icon"`
	grant_type		string		`json:"grant_type"`
	secret			string 		`json:"secret"`
	redirect_uris	[]string	`json:"redirect_uris"`
}

type DiscordAuthResponse struct {
	access_token	string		`json:"access_token"`
	token_type		string		`json:"token_type"`
	expires_in 		int			`json:"expires_in"`
	refresh_token	string 		`json:"refresh_token"`
	scope			string 		`json:"scope"`
}

func GetJsonBytes(id string, name string, description string, icon string, secret string, redirect_uris []string) ([]byte, error) {

	auth_request := DiscordAuthRequest{id, name, description,icon, "authorization_code",secret, redirect_uris}

	toRet, err := json.Marshal(auth_request)
	if (err != nil) {
		return nil, err
	}

	return toRet, nil
}

// https://discordapp.com/api/oauth2/token
func MakeDiscordAuthPOST(json []byte) {

	client := http.Client{}

	url_string := "https://discordapp.com/api/oauth2/token"

	data := url.Values{}
	data.Add("client_id", "462698059701420052")
	data.Add("client_secret", "z976rWnj8rgxiKjOD_rj578GcCS3hrd5")
	data.Add("grant_type", "authorization_code")
	data.Add("code", "")

	//request, err := http.NewRequest("POST", url_string, bytes.NewBuffer(json))
	request, err := http.NewRequest("POST", url_string, strings.NewReader(data.Encode()))
	if (err != nil) {
		panic(err)
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if (err != nil) {
		panic(err)
	}

	text, err := ioutil.ReadAll(response.Body)
	if (err != nil) {
		panic(err)
	}

	fmt.Printf("response: %s", text)
}