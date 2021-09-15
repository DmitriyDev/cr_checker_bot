package main

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type SyncGithubClient struct {
	ctx    context.Context
	token  oauth2.TokenSource
	Client *github.Client
}

func (sgc SyncGithubClient) New(oauthToken string) GithubClient {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: oauthToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &SyncGithubClient{
		ctx:    ctx,
		token:  ts,
		Client: client,
	}
}

func (sgc SyncGithubClient) GetFullInfo(info PRInfo) (PRResponseInfo, error) {
	prInfoList := map[int]PRInfo{}
	prInfoList[info.Id] = info
	prs, err := sgc.GetMultipleFullInfo(prInfoList)
	return *prs[info.Id], err
}

func (sgc SyncGithubClient) GetMultipleFullInfo(infos map[int]PRInfo) (map[int]*PRResponseInfo, error) {

	dataMap := map[int]*PRResponseInfo{}

	for _, info := range infos {
		prData := PRResponseInfo{Id: info.Id, Info: info, Error: []error{}}

		pRequest, errPR := sgc.GetPullRequestInfo(info)
		if errPR != nil {
			prData.Error = append(prData.Error, errPR)
			continue
		}
		pRequestReviews, errRV := sgc.GetReviews(info)
		if errRV != nil {
			prData.Error = append(prData.Error, errRV)
			continue
		}

		prData.PullRequest = pRequest
		prData.Reviews = pRequestReviews
		dataMap[info.Id] = &prData
	}

	return dataMap, nil
}

func (sgc SyncGithubClient) GetPullRequestInfo(info PRInfo) (github.PullRequest, error) {
	pr, _, err := sgc.Client.PullRequests.Get(sgc.ctx, info.Owner, info.Repo, info.Id)

	if err != nil {
		return github.PullRequest{}, err
	}
	return *pr, nil
}

func (sgc SyncGithubClient) GetReviews(info PRInfo) ([]*github.PullRequestReview, error) {
	opt := &github.ListOptions{}
	pr, _, err := sgc.Client.PullRequests.ListReviews(sgc.ctx, info.Owner, info.Repo, info.Id, opt)

	if err != nil {
		return []*github.PullRequestReview{}, err
	}
	return pr, nil
}
