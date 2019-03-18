package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tidwall/gjson"
)

func main() {

	token, err := signin()
	if err != nil {
		return
	}

	signout(token)
}

func signin() (string, error) {

	buf := bytes.NewBufferString(`{
  "login": "thewulf",
  "password": "Y3YGgGwj/W2&m84o;" }`)
	resp, err := http.Post("https://api.fleetyards.net/v1/sessions",
		"application/json", buf)
	if err != nil {
		log.Print("Error getting session ", err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("Error reading session body ", err)
		return "", err
	}

	token := gjson.Get(string(body), "token")

	return token.Str, nil
}

func signout(token string) {

	req, err := http.NewRequest("DELETE", "https://api.fleetyards.net/v1/sessions", nil)
	if err != nil {
		log.Fatal("Error creating delete request ", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	respd, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending delete ", err)
	}

	defer respd.Body.Close()
	body2, err := ioutil.ReadAll(respd.Body)
	if err != nil {
		log.Fatal("Error reading delete body ", err)
	}

	msg := gjson.Get(string(body2), "message")

	log.Print("delete msg ", msg.Str)

}
