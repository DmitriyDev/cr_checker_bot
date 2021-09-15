package main

const HelpOperation = "HELP"

type HelpHandler struct {
	commandType  string
	botProcessor BotProcessor
	response     CommandResponse
}

func (h HelpHandler) New(bp BotProcessor, cmdType string) CommandHandler {
	return &HelpHandler{botProcessor: bp, commandType: cmdType}
}

func (h *HelpHandler) Execute() (CommandResponse, error) {

	cr := CommandResponse{}.New()

	mes := "Available commands: \n\n" +
		"`!help` - Show available commands \n" +
		"`!add2team <slack user> <github login>` - (Admin only) add user to team\n" +
		"`!rm_user <slack user>` - (Admin only) remove user from team\n" +
		"`!team` - Show all users from team\n" +
		"`!stats` - Show all tracked PRs\n" +
		"`!stats_me` - Show all tracked PRs by user, who asked\n" +
		"`!update` - Update PR statuses and show new stats\n" +
		"`!invalidateCache` - (Admin only)Invalidate cached configs to reread from DB\n" +
		"`<PR url> [...<PR url>]` - Add PR to track\n"

	cr.SetReply(mes)
	cr.AddLog(HelpOperation, "Response")

	return cr, nil
}
