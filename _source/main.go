package main

import (
	"github.com/slack-go/slack/slackevents"
	"os"
)

var CredentialsManager Credentials

func main() {
	CredentialsManager = Credentials{}.Init()

	appToken := os.Getenv("SLACK_APP_TOKEN")
	botToken := os.Getenv("SLACK_BOT_TOKEN")

	sa := BotCommunicator{}.New(appToken, botToken)
	sa.Run()
}

func mentionedMessageHandler(sa *BotCommunicator, ev *slackevents.AppMentionEvent) {

	bp := BotProcessor{}.Init(sa, ev)
	cHandler := bp.GetCommandHandler()

	cmdRes, err := cHandler.Execute()

	logger := Logger{ev.Channel, ev.User}

	if err != nil {
		logger.log(LogMessage{}.Error("GENERAL_HANDLER", err.Error()))
		return
	}

	for _, lm := range cmdRes.LogMessage() {
		logger.log(lm)
	}

	for _, em := range cmdRes.ErrorMessage() {
		logger.log(em)
	}

	if cmdRes.doReply {
		bp.replyMessage(cmdRes.ReplyMessage())
	}
}
