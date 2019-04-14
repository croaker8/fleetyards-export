package fleetyards

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

// Signin - connect to fleetyards.net session and return authorization token
func Signin(user, pass string) (string, error) {

	// build credentials json
	buf := bytes.NewBufferString(fmt.Sprintf(`{
  "login": "%s",
  "password": "%s" }`, user, pass))

	// connect to fleetyards
	resp, err := http.Post("https://api.fleetyards.net/v1/sessions",
		"application/json", buf)
	if err != nil {
		fmt.Printf("Error connecting to fleetyards.net session: %s\n", err)
		return "", err
	}

	// read response body json and extract token field
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading fleetyards.net session connect body: %s\n", err)
		return "", err
	}
	token := gjson.Get(string(body), "token")

	return token.Str, nil
}

// Signout - disconnect from fleetyards session
func Signout(token string) error {

	// create request
	req, err := http.NewRequest("DELETE", "https://api.fleetyards.net/v1/sessions", nil)
	if err != nil {
		fmt.Printf("Error creating fleetyards.net session delete request: %s\n", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending fleetyards.net session delete request: %s\n", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {

		// read delete response message
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading fleetyards.net session delete response body: %s\n", err)
		}
		msg := gjson.Get(string(body), "message")
		fmt.Print("Session delete response from fleetyards.net: ", msg.Str)

	}

	return nil
}
