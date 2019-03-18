package main

import (
	"bytes"
	"fmt"
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

	getFleet(token)

	signout(token)
}

func getFleet(token string) {

	modList := make([]string, 0)

	for page := 1; ; page++ {

		url := fmt.Sprintf("https://api.fleetyards.net/v1/vehicles?page=%d", page)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Print("Error creating get vehicles request ", err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		req.Header.Add("Accept", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error sending get vehicles ", err)
		}

		//log.Print("Status code ", resp.StatusCode)

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading get vehicles body ", err)
		}

		bodyString := string(body)

		count := 0

		mods := gjson.Get(bodyString, "#.model.name")
		mods.ForEach(func(key, value gjson.Result) bool {
			modList = append(modList, value.String())
			count++
			return true // keep iterating
		})

		if count == 0 {
			break
		}
	}

	log.Printf("Got ship count %d", len(modList))
	for _, mod := range modList {
		log.Print(mod)
	}

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
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending delete ", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading delete body ", err)
	}

	msg := gjson.Get(string(body), "message")

	log.Print("delete msg ", msg.Str)

}
