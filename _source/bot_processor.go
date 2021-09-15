package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/slack-go/slack/slackevents"
	"os"
	"strings"
)

type BotProcessor struct {
	socketApp   *BotCommunicator
	event       *slackevents.AppMentionEvent
	credentials ChannelCredentials
}

var cachedDb = map[string]DB{}

func (bp BotProcessor) Init(sa *BotCommunicator, ev *slackevents.AppMentionEvent) BotProcessor {

	cc := CredentialsManager.Get(ev.Channel)

	return BotProcessor{
		socketApp:   sa,
		event:       ev,
		credentials: cc,
	}
}

func (bp *BotProcessor) GetDatabase() *DB {
	var db DB

	dbPath := os.Getenv("DB_DIR") + bp.credentials.DbName + ".csv"

	if _, ok := cachedDb[bp.credentials.DbName]; !ok {
		db = DB{}.Init(dbPath)
		cachedDb[dbPath] = db
	} else {
		db = cachedDb[dbPath]
	}

	return &db
}

func (bp *BotProcessor) GetGithubClient() GithubClient {

	switch os.Getenv("GITHUB_COMMUNICATION_MODE") {
	case "ASYNC":
		gc := AsyncGithubClient{}.New(bp.credentials.GithubToken)
		return gc
	case "SYNC":
		gc := SyncGithubClient{}.New(bp.credentials.GithubToken)
		return gc
	default:
		gc := SyncGithubClient{}.New(bp.credentials.GithubToken)
		return gc
	}

	return nil
}

func (bp BotProcessor) GetCommandHandler() CommandHandler {

	cMessage := bp.event.Text
	switch true {
	case strings.Contains(cMessage, "!help"):
		return HelpHandler{}.New(bp, HelpOperation)

	case strings.Contains(cMessage, "!team"):
		return TeamHandler{}.New(bp, TeamGetUsersOperation)
	case strings.Contains(cMessage, "!add2team"):
		return TeamHandler{}.New(bp, TeamAddUserOperation)
	case strings.Contains(cMessage, "!rm_user"):
		return TeamHandler{}.New(bp, TeamRemoveUserOperation)

	case strings.Contains(cMessage, "!invalidateCache"):
		return SystemHandler{}.New(bp, SystemInvalidateCacheOperation)

	case strings.Contains(cMessage, "!stats_me"):
		return PROperationHandler{}.New(bp, PROperationStatsMe)
	case strings.Contains(cMessage, "!stats"):
		return PROperationHandler{}.New(bp, PROperationStats)
	case strings.Contains(cMessage, "!update"):
		return PROperationHandler{}.New(bp, PROperationUpdate)
	default:
		return PROperationHandler{}.New(bp, PROperationUrlsRead)
	}
}

func (bp *BotProcessor) isUserAdmin() bool {
	adminUsers := strings.Split(os.Getenv("ADMIN_USERS"), ",")
	for _, au := range adminUsers {
		if au == bp.event.User {
			return true
		}
	}
	return false
}

func (bp *BotProcessor) replyMessage(message string) {
	bp.socketApp.ReplyMessage(bp.event.Channel, bp.event.TimeStamp, message)
}

func (bp *BotProcessor) formatPrInfo(pRequest github.PullRequest, prStatus PRStatus) string {
	status := ":no_entry_sign:"
	if prStatus.approves >= 2 {
		status = ":white_check_mark:"
	}

	team := bp.credentials.Team

	author := pRequest.GetUser().GetLogin()
	excludeLogins := []GithubLogin{GithubLogin(author)}
	for _, approver := range prStatus.approvedBy {
		ghL := GithubLogin(approver)
		excludeLogins = append(excludeLogins, ghL)
	}

	filteredTeam := team.NewTeamByGithubExcept(excludeLogins)

	var slakApprovers []string
	if bp.credentials.SilenceMode {
		for _, tu := range filteredTeam.Users {
			slakApprovers = append(slakApprovers, " `"+tu.GetGithubLogin()+"` ")
		}
	} else {
		for _, tu := range filteredTeam.Users {
			slakApprovers = append(slakApprovers, "<@"+tu.GetSlackId()+">")
		}
	}

	return fmt.Sprintf(
		"Status: %s\n"+
			"Title: %s\n"+
			"Author: %s\n"+
			"Url: %s\n"+
			"Review Comments: %d \n"+
			"Approves: %d\n"+
			"Approved By: %s\n"+
			"Waiting Approves from: %s"+
			"\n\n",
		status,
		pRequest.GetTitle(),
		pRequest.GetUser().GetLogin(),
		pRequest.GetHTMLURL(),
		prStatus.comments,
		prStatus.approves,
		strings.Join(prStatus.approvedBy, ", "),
		strings.Join(slakApprovers, ", "),
	)
}
