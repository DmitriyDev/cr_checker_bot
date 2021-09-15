package main

import "fmt"

type LogMessage struct {
	Type    string
	Action  string
	Message string
}

type Logger struct {
	Channel string
	User    string
}

func (lm LogMessage) General(action string, message string) LogMessage {
	return LogMessage{
		Type:    "General",
		Action:  action,
		Message: message,
	}
}

func (lm LogMessage) Error(action string, message string) LogMessage {
	return LogMessage{
		Type:    "Error",
		Action:  action,
		Message: message,
	}
}

func (l *Logger) log(lm LogMessage) {
	mes := fmt.Sprintf("[%s][%s][%s] (%s) %s", lm.Type, lm.Action, l.Channel, l.User, lm.Message)
	fmt.Println(mes)
}
