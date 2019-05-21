package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	fmt.Println("Init...")
	handleRequest()
}

func handleRequest() {
	http.HandleFunc("/bot", botFunc)
	http.HandleFunc("/status", statusPage)
	http.ListenAndServe(getPort(), nil)
}

func statusPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "It's ok!")
}

func botFunc(w http.ResponseWriter, req *http.Request) {
	if bot == nil {
		bot, _ = linebot.New(
			"ae6P1wQm9pDtBXz1TQNnAqWJSUHvIiUl0GWPJNvLK8MQxYuPIaqaP+Kea9H6QcnyVCyw2iJILvy00zXyV/B9nIB+NAeP9P9da7HZxbk0atcm2tYeuXngrKaMBMWwMy3msa5PEluN2bGu0JI7enTELwdB04t89/1O/w1cDnyilFU=",
			"097b8b8f50bbe56a02030db996b7b69d",
		)
	}

	events, err := bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
	fmt.Fprint(w, "bot bot")
}

func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("No Port In Heroku" + port)
	}
	return ":" + port
}
