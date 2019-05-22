package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var lastReq string

func main() {
	fmt.Println("Create File Main.go Success")
	handleRequest()
}

func handleRequest() {
	http.HandleFunc("/callback", botFunc)
	http.HandleFunc("/status", statusPage)
	http.HandleFunc("/lastrq", lastRequest)
	http.ListenAndServe(getPort(), nil)
}

func statusPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "It's ok!")
}

func lastRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, lastReq)
}

func botFunc(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	lastReq := string(bodyBytes)
	//print raw response body for debugging purposes
	fmt.Println("\n\n", lastReq, "\n\n")

	//reset the response body to the original unread state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	url := "https://api.line.me/v2/bot/message/reply"
	accToken := "ae6P1wQm9pDtBXz1TQNnAqWJSUHvIiUl0GWPJNvLK8MQxYuPIaqaP+Kea9H6QcnyVCyw2iJILvy00zXyV/B9nIB+NAeP9P9da7HZxbk0atcm2tYeuXngrKaMBMWwMy3msa5PEluN2bGu0JI7enTELwdB04t89/1O/w1cDnyilFU="
	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", " Bearer "+accToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(body)
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("No Port In Heroku" + port)
	}
	return ":" + port
}
