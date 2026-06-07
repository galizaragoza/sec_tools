package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gen2brain/beeep"
)

var (
	tgToken = os.Getenv("NOTIFY_CHAT_TOKEN")
	chatID  = os.Getenv("NOTIFY_CHAT_ID")
	chatURL = "https://api.telegram.org/bot" + tgToken + "/sendMessage"
	msg     []string
)

func tgNoti(msg string) (err error) {
	res, err := http.PostForm(chatURL,
		url.Values{
			"text":    {msg},
			"chat_id": {chatID},
		})
	if err != nil {
		return fmt.Errorf("error making the post request to the bot in tgNoti(): %w", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("post request returned code %d in tgNoti()", res.StatusCode)
	} else {
		return nil
	}
}

func osNoti(msg string) {
	err := beeep.Alert("Notify alert", msg, "bell.png")
	if err != nil {
		fmt.Errorf("error in osNoti(): %w", err)
	}
}

func main() {
	beeep.AppName = "Notify"

	if len(os.Args) >= 2 {
		msg = os.Args[1:]
	} else {
		err := fmt.Errorf("not enough arguments to parse the message on main()")
		fmt.Println(err)
	}

	joinedMsg := strings.Join(msg, " ")
	err := tgNoti(joinedMsg)
	if err != nil {
		fmt.Printf("Error in tgNoti func:", err)
	}

	osNoti(joinedMsg)
}
