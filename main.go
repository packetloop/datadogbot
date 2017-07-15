package main

import (
	"fmt"
	"log"
	"os"
	"unicode"

	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)

	slack.SetLogger(logger)
	api.SetDebug(true)

	for msg := range rtm.IncomingEvents {
		fmt.Printf("Event Received: %+#v\n", msg)
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			// fmt.Println("Infos:", ev.Info)
			// fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace #general with your Channel ID
			// rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#genera"))

		case *slack.MessageEvent:
			fmt.Printf("ALERT BOT ID: %+#v\n", ev.BotID)
			fmt.Printf("ALERT Title: %s\n", ev.Attachments[0].Title)
			fmt.Printf("ALERT Text : %s\n", stringMinifier(ev.Attachments[0].Text))
			fmt.Printf("ALERT Channel: %+#v\n", ev.Channel)
			// log to stdout syslog format with
			/*
				alert_type: triggered/recovered
				message: text
				alert_name: [atlasdossubmissionsbucket] [backup-artifacts] [failure] [staging]
			*/
			// Ignore
			// fmt.Printf("Message: %v\n", ev)
			// fmt.Printf("Subtype: %v\n", ev.SubType)
			// fmt.Printf("alert message: %+#v\n", ev.SubMessage)

		case *slack.RTMError:
			// fmt.Printf("Error: %s\n", ev.Error())

		default:
			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

// Minify slack event text
func stringMinifier(in string) (out string) {
	white := false
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				out = out + " "
			}
			white = true
		} else {
			out = out + string(c)
			white = false
		}
	}
	return
}
