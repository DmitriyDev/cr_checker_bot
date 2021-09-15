package main

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
	"strings"
)

type BotCommunicator struct {
	api    *slack.Client
	client *socketmode.Client
}

func (bc BotCommunicator) New(aToken string, bToken string) BotCommunicator {
	bc.validate(aToken, bToken)

	ac := slack.New(
		bToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(aToken),
	)

	sc := socketmode.New(
		ac,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	return BotCommunicator{api: ac, client: sc}

}

func (bc BotCommunicator) PostMessage(channelId string, message string) {
	_, _, err := bc.api.PostMessage(channelId,
		slack.MsgOptionText(message, false),
		slack.MsgOptionAsUser(false),
		slack.MsgOptionIconEmoji(":male-detective:"),
	)
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
	}
}
func (bc BotCommunicator) ReplyMessage(channelId string, ts string, message string) {
	_, _, err := bc.api.PostMessage(channelId,
		slack.MsgOptionText(message, false),
		slack.MsgOptionTS(ts),
		slack.MsgOptionAsUser(false),
		slack.MsgOptionUsername("PullRequest Checker"),
		slack.MsgOptionIconEmoji(":male-detective:"),
	)
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
	}
}

func (bc BotCommunicator) validate(aToken string, bToken string) {

	if !strings.HasPrefix(aToken, "xapp-") {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must have the prefix \"xapp-\".")
		os.Exit(1)
	}

	if !strings.HasPrefix(bToken, "xoxb-") {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}
}

func (bc *BotCommunicator) Run() {
	go func() {
		for evt := range bc.client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)
					continue
				}

				fmt.Printf("Event received: %+v\n", eventsAPIEvent)

				bc.client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						mentionedMessageHandler(bc, ev)
					case *slackevents.MemberJoinedChannelEvent:
						fmt.Printf("user %q joined to channel %q", ev.User, ev.Channel)
					}
				default:
					bc.client.Debugf("unsupported Events API event received")
				}
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	bc.client.Run()
}
