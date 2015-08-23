package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins/echo"
)

func main() {
	var token string
	flag.StringVar(&token, "token", os.Getenv("SLACK_BOT_TOKEN"), "SlackのBotToken")
	flag.Parse()

	bot, err := slackbot.NewBotContext(token)
	if err != nil {
		panic(err)
	}
	bot.AddPlugin("echo", echo.Plugin{})

	bot.WebSocketRTM()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(":8000", nil)
}
