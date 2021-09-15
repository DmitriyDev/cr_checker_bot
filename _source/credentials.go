package main

import (
	"encoding/json"
	"os"
)

type Credentials struct {
	channelsCred map[string]ChannelCredentials
	dirPath      string
}
type ChannelCredentials struct {
	DbName      string `json:"db_name"`
	GithubToken string `json:"github_token"`
	SilenceMode bool   `json:"silence_mode"`
	Team        Team   `json:"team"`
}

func (c Credentials) Init() Credentials {
	return Credentials{
		channelsCred: map[string]ChannelCredentials{},
		dirPath:      os.Getenv("CREDENTIALS_DIR"),
	}
}

func (c *Credentials) getChannelStorage(channel string) Storage {
	return Storage{}.New(c.dirPath + channel + ".csv")
}

func (c *Credentials) defaultCredentials(channel string) ChannelCredentials {
	return ChannelCredentials{
		DbName:      channel,
		GithubToken: os.Getenv("GITHUB_TOKEN"),
		SilenceMode: os.Getenv("SILENT_MODE") == "true",
		Team:        Team{}.New(),
	}
}

func (c *Credentials) getFromStorage(channel string) ChannelCredentials {

	storage := c.getChannelStorage(channel)

	var cc ChannelCredentials
	if storage.exists() {
		rawCredentials := storage.Get()
		if rawCredentials != "" {
			cc = ChannelCredentials{}
			_ = json.Unmarshal([]byte(rawCredentials), &cc)
			return cc
		}
	}

	cc = c.defaultCredentials(channel)
	ccRaw, _ := json.Marshal(cc)
	storage.ReplaceWith(string(ccRaw))

	return cc
}

func (c *Credentials) dump(channel string) {
	storage := c.getChannelStorage(channel)

	ccRaw, _ := json.Marshal(c.channelsCred[channel])
	storage.ReplaceWith(string(ccRaw))
}

func (c *Credentials) Get(channel string) ChannelCredentials {
	if _, ok := c.channelsCred[channel]; !ok {
		c.channelsCred[channel] = c.getFromStorage(channel)
	}
	return c.channelsCred[channel]
}

func (c *Credentials) Update(channel string, cc ChannelCredentials) {
	c.channelsCred[channel] = cc
	c.dump(channel)
}

func (c *Credentials) InvalidateCache() bool {
	c.channelsCred = map[string]ChannelCredentials{}
	return true
}
