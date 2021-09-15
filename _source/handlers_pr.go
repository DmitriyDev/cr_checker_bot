package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const PROperationStats = "STATS"
const PROperationStatsMe = "STATS_ME"
const PROperationUpdate = "UPDATE"
const PROperationUrlsRead = "URLS_READ"

type PROperationHandler struct {
	commandType  string
	botProcessor BotProcessor
	response     CommandResponse
}

func (ph PROperationHandler) New(bp BotProcessor, cmdType string) CommandHandler {
	return &PROperationHandler{botProcessor: bp, commandType: cmdType}
}

func (ph *PROperationHandler) Execute() (CommandResponse, error) {

	switch ph.commandType {
	case PROperationStats:
		return ph.HandleStats()
	case PROperationStatsMe:
		return ph.HandleStatsMe()
	case PROperationUpdate:
		return ph.HandleUpdate()
	case PROperationUrlsRead:
		return ph.HandleUrlRead()
	default:
		return CommandResponse{}, errors.New("invalid operation")
	}
}

func (ph *PROperationHandler) HandleStats() (CommandResponse, error) {

	cr := CommandResponse{}.New()

	message := "No Pull requests available to show"
	reply := ""
	rcc := ph.botProcessor.GetDatabase().ReadRecords()

	for _, rec := range rcc {
		prStatus := extractGHPRInfo(rec.ReviewData)
		reply += ph.botProcessor.formatPrInfo(rec.Data, prStatus)
	}

	if reply != "" {
		message = "Pull requests Status chart\n\n" + reply
	}

	cr.AddLog(PROperationStats, "Response")
	cr.SetReply(message)

	return cr, nil
}

func (ph *PROperationHandler) HandleStatsMe() (CommandResponse, error) {

	cr := CommandResponse{}.New()

	message := "No Pull requests available to show"
	reply := ""

	creds := CredentialsManager.Get(ph.botProcessor.event.Channel)
	u, err := creds.Team.getBySlack(SlackId(ph.botProcessor.event.User))

	if err == nil {
		rcc := ph.botProcessor.GetDatabase().ReadRecords()

		for _, rec := range rcc {
			if rec.Data.GetUser().GetLogin() == u.GetGithubLogin() {
				prStatus := extractGHPRInfo(rec.ReviewData)
				reply += ph.botProcessor.formatPrInfo(rec.Data, prStatus)
			}
		}

		if reply != "" {
			message = "Pull requests status chart by " + u.GetGithubLogin() + "\n\n" + reply
		}

	} else {
		message = "You not added to the team"

	}

	cr.AddLog(PROperationStatsMe, "Response")
	cr.SetReply(message)

	return cr, nil
}

func (ph *PROperationHandler) HandleUpdate() (CommandResponse, error) {

	cr := CommandResponse{}.New()

	githubClient := ph.botProcessor.GetGithubClient()

	reply := ""
	rcc := ph.botProcessor.GetDatabase().ReadRecords()
	removed := 0

	prInfoList := map[int]PRInfo{}
	for _, rec := range rcc {
		prInfo := rec.BaseInfo
		prInfoList[prInfo.Id] = prInfo
	}

	prRSMap, err := githubClient.GetMultipleFullInfo(prInfoList)

	if err != nil {
		cr.AddError(PROperationUrlsRead, err)
	}

	for _, prRs := range prRSMap {
		if len(prRs.Error) != 0 {
			for _, err := range prRs.Error {
				cr.AddError(PROperationUrlsRead, err)
			}
			continue
		}
		prStatus := extractGHPRInfo(prRs.Reviews)
		reply += ph.botProcessor.formatPrInfo(prRs.PullRequest, prStatus)

		if prStatus.approves >= 2 {
			removed += 1
			ph.botProcessor.GetDatabase().Remove(prRs.PullRequest.GetID())
			continue
		}
		ph.botProcessor.GetDatabase().Add(prRs.Info, prRs.PullRequest, prRs.Reviews)
	}

	message := ""
	if removed > 0 {
		message = fmt.Sprintf("DB Updated (%d removed).", removed)
	}

	if reply != "" {
		message = fmt.Sprintf("%s New Statuses:\n\n %s", message, reply)

	} else {
		message = fmt.Sprintf("%s Nothing left\n", message)
	}

	cr.AddLog(PROperationUpdate, "Response")
	cr.SetReply(message)

	return cr, nil

}

func (ph *PROperationHandler) HandleUrlRead() (CommandResponse, error) {

	message := "No links found"
	cr := CommandResponse{}.New()

	txt, perr := strconv.Unquote(`"` + ph.botProcessor.event.Text + `"`)
	if perr != nil {
		txt = ph.botProcessor.event.Text
	}

	re := regexp.MustCompile("(?m)(https://github.com/\\w*/\\w*/pull/\\d*)")
	urls := re.FindAllString(txt, -1)

	if len(urls) == 0 {
		cr.SetReply("Invalid command!\nCall `!help` to receive commands list")
		return cr, nil
	}

	reply := ""

	githubClient := ph.botProcessor.GetGithubClient()

	prInfoList := map[int]PRInfo{}
	for i := range urls {
		prInfo, _ := getPrInfo(urls[i])
		prInfoList[prInfo.Id] = prInfo
	}

	prRSMap, err := githubClient.GetMultipleFullInfo(prInfoList)

	if err != nil {
		cr.AddError(PROperationUrlsRead, err)
	}

	for _, prRs := range prRSMap {
		if len(prRs.Error) != 0 {
			for _, err := range prRs.Error {
				cr.AddError(PROperationUrlsRead, err)
			}
			continue
		}
		prStatus := extractGHPRInfo(prRs.Reviews)
		reply += ph.botProcessor.formatPrInfo(prRs.PullRequest, prStatus)
		ph.botProcessor.GetDatabase().Add(prRs.Info, prRs.PullRequest, prRs.Reviews)
	}

	if reply != "" {
		message = "Added PR links\n\n" + reply
	}

	cr.AddLog(PROperationUrlsRead, "Response")
	cr.SetReply(message)

	return cr, nil
}
