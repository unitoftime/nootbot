package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// utility functions that used by multiple cmds

func GetJson(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		fmt.Println("error with GetJson")
		return err
	}
	defer r.Body.Close()

	// b, _ := io.ReadAll(r.Body)
	// fmt.Println(string(b))

	return json.NewDecoder(r.Body).Decode(target)
}

func ReadFile(URL string) ([]byte, error) {
	//Get the response bytes from the url
	resp, err := http.Get(URL)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return []byte{}, errors.New("Received non 200 response code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.New("Read all failed")
	}

	//fmt.Println(string(body))

	return body, nil
}
