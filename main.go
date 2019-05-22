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
	var lineEvents = LineEvents{}
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	lastReq = string(bodyBytes)
	//print raw response body for debugging purposes
	fmt.Println("\n\n", lastReq, "\n\n")

	//reset the response body to the original unread state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(r.Body).Decode(&lineEvents); err != nil {
		return
	}

	for _, lineEvent := range lineEvents.Events {
		if lineEvent.Type == "message" {
			reply := &Reply{ReplyToken: lineEvent.ReplyToken, Message: Message{Type: lineEvent.Type, Text: lineEvent.Message.Text}}
			fmt.Print(reply)

			url := "https://api.line.me/v2/bot/message/reply"
			accToken := "ae6P1wQm9pDtBXz1TQNnAqWJSUHvIiUl0GWPJNvLK8MQxYuPIaqaP+Kea9H6QcnyVCyw2iJILvy00zXyV/B9nIB+NAeP9P9da7HZxbk0atcm2tYeuXngrKaMBMWwMy3msa5PEluN2bGu0JI7enTELwdB04t89/1O/w1cDnyilFU="

			jsonStr, _ := json.Marshal(*reply)

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", " Bearer "+accToken)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			lastReq = string(jsonStr)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, "It's oks!")
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("No Port In Heroku" + port)
	}
	return ":" + port
}

type LineEvents struct {
	Events      []Events `json:"events"`
	Destination string   `json:"destination"`
}
type Source struct {
	UserID  string `json:"userId"`
	GroupID string `json:"groupId"`
	Type    string `json:"type"`
}
type Message struct {
	Type string `json:"type"`
	ID   string `json:"id,omitempty"`
	Text string `json:"text"`
}
type Events struct {
	Type       string  `json:"type"`
	ReplyToken string  `json:"replyToken"`
	Source     Source  `json:"source"`
	Timestamp  int64   `json:"timestamp"`
	Message    Message `json:"message"`
}

type Reply struct {
	ReplyToken string  `json:"replyToken"`
	Message    Message `json:"messages"`
}
