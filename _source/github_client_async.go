package main

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"sync"
)

type AsyncGithubClient struct {
	ctx    context.Context
	token  oauth2.TokenSource
	Client *github.Client
}

func (agc AsyncGithubClient) New(oauthToken string) GithubClient {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: oauthToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &AsyncGithubClient{
		ctx:    ctx,
		token:  ts,
		Client: client,
	}
}

func (agc AsyncGithubClient) GetFullInfo(info PRInfo) (PRResponseInfo, error) {
	prInfoList := map[int]PRInfo{}
	prInfoList[info.Id] = info
	prs, err := agc.GetMultipleFullInfo(prInfoList)
	return *prs[info.Id], err
}

func (agc AsyncGithubClient) GetMultipleFullInfo(infos map[int]PRInfo) (map[int]*PRResponseInfo, error) {
	var wg sync.WaitGroup

	dataMap := map[int]*PRResponseInfo{}
	for _, info := range infos {
		prData := PRResponseInfo{Id: info.Id, Info: info, Error: []error{}}
		dataMap[info.Id] = &prData

		wg.Add(2)
		go agc.asyncPrInfo(&wg, dataMap, info.Id)
		go agc.asyncPrReviews(&wg, dataMap, info.Id)

	}

	wg.Wait()

	return dataMap, nil
}

func (agc AsyncGithubClient) asyncPrInfo(wg *sync.WaitGroup, r map[int]*PRResponseInfo, id int) {
	defer wg.Done()
	info := r[id].Info
	pr, _, err := agc.Client.PullRequests.Get(agc.ctx, info.Owner, info.Repo, info.Id)
	if err != nil {
		r[id].Error = append(r[id].Error, err)
	} else {
		r[id].PullRequest = *pr
	}
}

func (agc AsyncGithubClient) asyncPrReviews(wg *sync.WaitGroup, r map[int]*PRResponseInfo, id int) {
	defer wg.Done()
	info := r[id].Info
	opt := &github.ListOptions{}
	rv, _, err := agc.Client.PullRequests.ListReviews(agc.ctx, info.Owner, info.Repo, info.Id, opt)
	if err != nil {
		r[id].Error = append(r[id].Error, err)
	} else {
		r[id].Reviews = rv
	}
}
