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
	BotName     string
	Text        string
	Environment string
}

var (
	botName, token string
	logger         *log.Logger
	rtm            *slack.RTM
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token = os.Getenv("SLACK_TOKEN")
	botName = os.Getenv("BOTNAME")

	logger, err = syslog.NewLogger(syslog.LOG_LOCAL3, log.Lmicroseconds)
	if err != nil {
		fmt.Println("Cannot set syslog")
		os.Exit(1)
	}

	// Explicitly add a trailing space. Set prefix does not add a trailing
	// space.
	logger.SetPrefix("slackbot ")
	slack.SetLogger(logger)
	api := slack.New(token)

	rtm = api.NewRTM()
}

func main() {
	go rtm.ManageConnection()

	event := make(chan alertEvent)
	data := make(chan *slack.MessageEvent)
	for i := 0; i <= 10; i++ {
		go parseAlert(data, event)
		go listen(event, rtm)
	}
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			// Ignore connected events

		case *slack.MessageEvent:
			data <- ev
		case *slack.RTMError:
			logger.Printf("rtm error: %s\n", ev.Error())

		default:
			// Ignore other events..
		}
	}
}

func listen(event chan alertEvent, rtm *slack.RTM) {
	for alert := range event {
		logger.Printf("%s %s %s %s %s description %s",
			alert.Channel, alert.BotName, alert.Environment, alert.Type,
			alert.Title, stringMinifier(alert.Text))
	}
}

func getBotName(botID string) string {
	evBot, err := rtm.GetBotInfo(botID)
	if err != nil {
		logger.Printf("getting bot name error %s", err.Error())
	}
	return evBot.Name
}

func getChannelName(channel string) string {
	ch, err := rtm.GetChannelInfo(channel)
	if err != nil {
		logger.Printf("getting channel name error %s", err.Error())
	}
	return ch.Name
}

// Parse alert
func parseAlert(data <-chan *slack.MessageEvent, alert chan<- alertEvent) {
	for ev := range data {
		// We only listen to alerts that comes from Datadog and ignore
		// all other alerts.
		evBotName := getBotName(ev.BotID)
		if evBotName == botName {
			channelName := getChannelName(ev.Channel)
			if evBotName == botName {
				event := strings.Split(ev.Attachments[0].Title, " ")
				var alertName []string
				for i, v := range event {
					if i == 0 || i == len(event)-1 {
						continue
					}
					alertName = append(alertName, alertTitleMinify(v))
				}

				alert <- alertEvent{
					Type:        strings.Trim(event[0], ": "),
					Channel:     channelName,
					BotName:     evBotName,
					Text:        ev.Attachments[0].Text,
					Environment: event[len(event)-1],
					Title:       strings.Join(alertName, " "),
				}
			}
		}
	}
}

// Minify alert title
func alertTitleMinify(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "["), "]")
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
