package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const TeamGetUsersOperation = "GET_USERS"
const TeamRemoveUserOperation = "RM_USER"
const TeamAddUserOperation = "ADD_USER"

type TeamHandler struct {
	commandType  string
	botProcessor BotProcessor
	response     CommandResponse
}

func (th TeamHandler) New(bp BotProcessor, cmdType string) CommandHandler {
	return &TeamHandler{botProcessor: bp, commandType: cmdType}
}

func (th *TeamHandler) Execute() (CommandResponse, error) {

	switch th.commandType {
	case TeamGetUsersOperation:
		return th.HandleTeamUsers()
	case TeamAddUserOperation:
		return th.HandleAddUserToTeam()
	case TeamRemoveUserOperation:
		return th.HandleRemoveUserFromTeam()
	default:
		return CommandResponse{}, errors.New("invalid operation")
	}
}

func (th *TeamHandler) HandleTeamUsers() (CommandResponse, error) {

	message := "No users added"
	cr := CommandResponse{}.New()

	creds := CredentialsManager.Get(th.botProcessor.event.Channel)

	mes := ""
	for _, tUser := range creds.Team.Users {
		mes += fmt.Sprintf("Slack: %s\n Github: %s\n\n", tUser.SlackId, tUser.GithubLogin)
	}

	if mes != "" {
		message = "Team:\n" + mes
	}

	cr.SetReply(message)
	cr.AddLog(TeamGetUsersOperation, "Requested")

	return cr, nil
}

func (th *TeamHandler) HandleRemoveUserFromTeam() (CommandResponse, error) {

	cr := CommandResponse{}.New()

	if !th.botProcessor.isUserAdmin() {
		cr.SetReply("No Access")
		cr.AddLog(TeamGetUsersOperation, "No Access")
		return cr, nil
	}

	mes := strings.Replace(th.botProcessor.event.Text, " ", " ", -1)
	re := regexp.MustCompile(`.*!rm_user\s*<@(\S*)>`)
	cmdParts := re.FindAllStringSubmatch(mes, -1)

	if len(cmdParts) == 0 {
		cr.SetReply("Invalid command! Expected: !rm_user SlackUser")
		return cr, nil
	}

	slackId := SlackId(cmdParts[0][1])

	creds := CredentialsManager.Get(th.botProcessor.event.Channel)

	user, err := creds.Team.getBySlack(slackId)

	if err != nil {
		cr.SetReply("User not found")
		cr.AddError(TeamRemoveUserOperation, err)
		return cr, nil
	}

	creds.Team.remove(user)
	CredentialsManager.Update(th.botProcessor.event.Channel, creds)

	cr.SetReply(fmt.Sprintf("User removed from team (slack: %s)", slackId))
	cr.AddLog(TeamRemoveUserOperation, fmt.Sprintf("User removed from team %s", slackId))
	return cr, nil
}

func (th *TeamHandler) HandleAddUserToTeam() (CommandResponse, error) {

	cr := CommandResponse{}.New()

	if !th.botProcessor.isUserAdmin() {
		cr.SetReply("No Access")
		cr.AddLog(TeamAddUserOperation, "No Access")
		return cr, nil
	}

	mes := strings.Replace(th.botProcessor.event.Text, " ", " ", -1)
	re := regexp.MustCompile(`.*!add2team\s*<@(\S*)>\s*(\S*)`)
	cmdParts := re.FindAllStringSubmatch(mes, -1)

	if len(cmdParts) == 0 {
		cr.SetReply("Invalid command! Expected: !rm_user SlackUser")
		return cr, nil
	}

	slackId := SlackId(cmdParts[0][1])
	githubLogin := GithubLogin(cmdParts[0][2])

	creds := CredentialsManager.Get(th.botProcessor.event.Channel)

	user, err := creds.Team.getBySlack(slackId)

	if err == nil {
		user.SlackId = slackId
		user.GithubLogin = githubLogin
	} else {
		user = creds.Team.createNewUser(slackId, githubLogin)
	}
	creds.Team.add(user)

	CredentialsManager.Update(th.botProcessor.event.Channel, creds)

	cr.SetReply(fmt.Sprintf("Add user to team:\nSlack: %s\nGithub: %s", slackId, githubLogin))
	cr.AddLog(TeamAddUserOperation, fmt.Sprintf("User added to team %s (%s)", slackId, githubLogin))
	return cr, nil
}
