package main

import "errors"

const SystemInvalidateCacheOperation = "INVALIDATE_CACHE"

type SystemHandler struct {
	commandType  string
	botProcessor BotProcessor
	response     CommandResponse
}

func (sh SystemHandler) New(bp BotProcessor, cmdType string) CommandHandler {
	return &SystemHandler{botProcessor: bp, commandType: cmdType}
}

func (sh *SystemHandler) Execute() (CommandResponse, error) {

	switch sh.commandType {
	case SystemInvalidateCacheOperation:
		return sh.InvalidateCache()
	default:
		return CommandResponse{}, errors.New("invalid operation")
	}
}

func (sh *SystemHandler) InvalidateCache() (CommandResponse, error) {

	messsage := "No access"
	cr := CommandResponse{}.New()

	cr.AddLog(SystemInvalidateCacheOperation, "Requested")

	if sh.botProcessor.isUserAdmin() {
		CredentialsManager.InvalidateCache()
		messsage = "Credentials cache invalidated"
		cr.AddLog(SystemInvalidateCacheOperation, "Credentials cache invalidated")
	}

	cr.SetReply(messsage)

	return cr, nil
}
