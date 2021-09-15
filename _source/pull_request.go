package main

import (
	"errors"
	"github.com/google/go-github/github"
	"regexp"
	"strconv"
)

type PRInfo struct {
	Owner string
	Repo  string
	Id    int
}

type PRStatus struct {
	approves   int
	comments   int
	approvedBy []string
}

func extractGHPRInfo(prvs []*github.PullRequestReview) PRStatus {

	st := PRStatus{approves: 0, comments: 0, approvedBy: []string{}}

	for _, rw := range prvs {
		if *rw.State == "COMMENTED" {
			st.comments += 1
			continue
		}
		if *rw.State == "APPROVED" {
			st.approves += 1
			st.approvedBy = append(st.approvedBy, rw.GetUser().GetLogin())
			continue
		}
	}

	return st
}

func getPrInfo(url string) (PRInfo, error) {

	re := regexp.MustCompile("https://github.com/(.*)/(.*)/pull/(\\d*)")
	res := re.FindAllStringSubmatch(url, 1)
	for i := range res {

		owner := res[i][1]
		repo := res[i][2]
		id, _ := strconv.Atoi(res[i][3])

		return PRInfo{
			Owner: owner,
			Repo:  repo,
			Id:    id,
		}, nil
	}

	return PRInfo{}, errors.New("No url found")

}
