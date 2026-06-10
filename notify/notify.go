package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
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
	cmd := exec.Command("/usr/bin/notify-send", "-a", "Notify", "-u", "normal", "-t", "2000", "-n", "bell.png",
		msg)
	err := cmd.Run()
	if err != nil {
		fmt.Errorf("error running os notification command on osNoti(): %w", err)
	}
}

func main() {
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
