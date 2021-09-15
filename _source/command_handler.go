package main

type CommandHandler interface {
	New(ah BotProcessor, cmdType string) CommandHandler
	Execute() (CommandResponse, error)
}

type CommandResponse struct {
	doReply      bool
	errors       []LogMessage
	logMessage   []LogMessage
	replyMessage string
}

func (cr CommandResponse) New() CommandResponse {
	return CommandResponse{
		true,
		[]LogMessage{},
		[]LogMessage{},
		"",
	}
}
func (cr *CommandResponse) AddError(actionType string, err error) {
	cr.errors = append(cr.errors, LogMessage{}.Error(actionType, err.Error()))
}

func (cr *CommandResponse) SetReply(message string) {
	cr.replyMessage = message
}

func (cr *CommandResponse) AddLog(actionType string, message string) {
	cr.errors = append(cr.errors, LogMessage{}.General(actionType, message))
}

func (cr *CommandResponse) ErrorMessage() []LogMessage {
	return cr.errors
}

func (cr *CommandResponse) LogMessage() []LogMessage {
	return cr.logMessage
}

func (cr *CommandResponse) ReplyMessage() string {
	return cr.replyMessage
}
