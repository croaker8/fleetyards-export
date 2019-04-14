package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

func main() {

	token, err := signin()
	if err != nil {
		return
	}

	fieldList, err := getFieldList()
	if err != nil {
		os.Exit(1)
	}

	if len(fieldList) == 0 {
		log.Fatal("No fields selected")
		os.Exit(2)
	}

	processFleet(token, "output.csv", fieldList)

	signout(token)
}

func processFleet(token, outFile string, fieldList []string) {

	fieldCount := len(fieldList)

	os.Remove(outFile)
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatal("Error creating output file ", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	// write csv header
	w.Write(fieldList)

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
		if bodyString == "[]" {
			break
		}

		mods := gjson.Parse(bodyString)

		mods.ForEach(func(key, value gjson.Result) bool {

			vals := make([]string, 0, fieldCount)

			log.Printf(value.String())

			for _, field := range fieldList {
				v := value.Get(field)
				vals = append(vals, v.String())
			}

			w.Write(vals)

			return true // keep iterating
		})

	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal("Error writing csv file ", err)
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

func getFieldList() ([]string, error) {

	list := make([]string, 0, 50)

	file, err := os.Open("field_list")
	if err != nil {
		fmt.Println("Error opening field_list file: ", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0:1] != "#" {
			list = append(list, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning lines from field list: ", err)
		return nil, err
	}

	return list, nil
}
