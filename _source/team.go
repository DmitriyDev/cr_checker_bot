package main

import (
	"errors"
)

type SlackId string
type GithubLogin string

type TeamUser struct {
	Id          int         `json:"id"`
	SlackId     SlackId     `json:"slack"`
	GithubLogin GithubLogin `json:"github"`
}

func (tu *TeamUser) GetSlackId() string {
	return string(tu.SlackId)
}

func (tu *TeamUser) GetGithubLogin() string {
	return string(tu.GithubLogin)
}

type Team struct {
	Users map[int]TeamUser `json:"team_user"`
}

func (t Team) New() Team {
	return Team{Users: map[int]TeamUser{}}
}

func (t *Team) newId() int {
	return len(t.Users) + 1
}

func (t *Team) createNewUser(su SlackId, gl GithubLogin) TeamUser {
	id := t.newId()
	return TeamUser{id, su, gl}
}

func (t *Team) add(user TeamUser) {
	t.Users[user.Id] = user
}

func (t *Team) remove(user TeamUser) {
	delete(t.Users, user.Id)
}

func (t *Team) getByGithub(gl GithubLogin) (TeamUser, error) {
	for _, tu := range t.Users {
		if gl == tu.GithubLogin {
			return tu, nil
		}
	}
	return TeamUser{}, errors.New("No slack user found")
}

func (t *Team) getBySlack(si SlackId) (TeamUser, error) {
	for _, tu := range t.Users {
		if si == tu.SlackId {
			return tu, nil
		}
	}
	return TeamUser{}, errors.New("No github user found")
}

func (t *Team) NewTeamBySlackExcept(sids []SlackId) Team {

	team := Team{Users: map[int]TeamUser{}}
	for _, tu := range t.Users {
		idFound := false

		for _, exSid := range sids {
			if tu.SlackId == exSid {
				idFound = true
				break
			}
		}

		if idFound {
			continue
		}
		team.add(tu)
	}
	return team
}


func (t *Team) NewTeamByGithubExcept(gids []GithubLogin) Team {

	team := Team{Users: map[int]TeamUser{}}
	for _, tu := range t.Users {
		idFound := false

		for _, exGid := range gids {
			if tu.GithubLogin == exGid {
				idFound = true
				break
			}
		}

		if idFound {
			continue
		}
		team.add(tu)
	}
	return team
}