package fleetyards

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

var defaultPageSize = 30

// ExportHangerToCsv - export fleetyards.net hanger data to CSV file
func ExportHangerToCsv(token, outFile string, fieldList []string) error {

	// create the new output file
	f, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("Error creating output file '%s': %s\n", outFile, err)
		return err
	}
	defer f.Close()

	// get csv writer and write header with field names
	w := csv.NewWriter(f)
	err = w.Write(fieldList)
	if err != nil {
		fmt.Printf("Error writing header line to output file: %s\n", err)
		return err
	}

	// loop through paged responses from fleetyards.net
	for page := 1; ; page++ {

		// get page of data
		pageList, err := getHangerPage(fieldList, page, token)
		if err != nil {
			// handing page failed
			return err
		}

		// list will be nil if there is no more data
		if pageList == nil {
			break
		}

		// write values to output file
		for _, vals := range pageList {
			err = w.Write(vals)
			if err != nil {
				fmt.Printf("Error writing row of data to output: %s\n", err)
				return err
			}
		}

	}

	// flush output data
	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Printf("Error flushing data to output file: %s\n", err)
		return err
	}

	return nil
}

func getHangerPage(fieldList []string, page int, token string) ([][]string, error) {

	// build request for the page
	url := fmt.Sprintf("https://api.fleetyards.net/v1/vehicles?page=%d", page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating get vehicles request for page %d: %s\n", page, err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")

	// request the page
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error getting vehicles for page %d: %s\n", page, err)
		return nil, err
	}

	// read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading page %d vehicles body: %s\n", page, err)
		return nil, err
	}

	// get string of body and return nil slice if no more data
	bodyString := string(body)
	if bodyString == "[]" {
		return nil, nil
	}

	// get count of fields for initializing vals slice
	fieldCount := len(fieldList)

	// parse the body
	mods := gjson.Parse(bodyString)

	// create list to return
	pageList := make([][]string, 0, defaultPageSize)

	// loop through each item in the page
	mods.ForEach(func(key, value gjson.Result) bool {

		// make slice to hold the field values
		vals := make([]string, 0, fieldCount)

		//fmt.Printf(value.String())

		// get the fields
		for _, field := range fieldList {
			v := value.Get(field)
			vals = append(vals, v.String())
		}

		// add field values to page list
		pageList = append(pageList, vals)

		return true // keep iterating
	})

	return pageList, nil

}
