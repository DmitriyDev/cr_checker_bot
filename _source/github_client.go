package main

import "github.com/google/go-github/github"

type GithubClient interface {
	New(oauthToken string) GithubClient
	GetFullInfo(info PRInfo) (PRResponseInfo, error)
	GetMultipleFullInfo(infos map[int]PRInfo) (map[int]*PRResponseInfo, error)
}

type PRResponseInfo struct {
	Id          int
	Info        PRInfo
	PullRequest github.PullRequest
	Reviews     []*github.PullRequestReview
	Error       []error
}
