package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"strings"
	"unicode"

	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
)

type alertEvent struct {
	Type        string
	Channel     string
	Title       string
	BotID       string
	Text        string
	Environment string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("SLACK_TOKEN")
	logger, err := syslog.NewLogger(syslog.LOG_LOCAL3, log.Lmicroseconds)
	logger.SetPrefix("slackbot")
	slack.SetLogger(logger)
	api := slack.New(token)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Printf("Event Received: %+#v\n", msg)
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			// Ignore connected events

		case *slack.MessageEvent:
			// We only listen to alerts that comes from Datadog Bot ID and ignore
			// all other alerts.
			if ev.BotID == "B1V1KFA0K" {
				alert := parseAlert(ev)
				channel, _ := rtm.GetChannelInfo(alert.Channel)
				if err != nil {
					logger.Printf("getting channel name error %s", err.Error())
				}
				bot, err := rtm.GetBotInfo(alert.BotID)
				if err != nil {
					logger.Printf("getting bot name error %s", err.Error())
				}
				fmt.Printf("channel string: %s\n", channel.Name)
				logger.Printf("channel %s bot %s environment %s alert type %s %s %s",
					channel.Name, bot.Name, alert.Environment, alert.Type,
					alert.Title, stringMinifier(alert.Text))
			}

		case *slack.RTMError:
			logger.Printf("rtm error: %s\n", ev.Error())

		default:
			// Ignore other events..
		}
	}
}

// Parse alert
func parseAlert(ev *slack.MessageEvent) alertEvent {
	event := strings.Split(ev.Attachments[0].Title, " ")
	var alertName []string
	for i, v := range event {
		if i == 0 || i == len(event)-1 {
			continue
		}
		alertName = append(alertName, v)
	}

	return alertEvent{
		Type:        strings.Trim(event[0], ": "),
		Channel:     ev.Channel,
		BotID:       ev.BotID,
		Text:        ev.Attachments[0].Text,
		Environment: event[len(event)-1],
		Title:       strings.Join(alertName, " "),
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
