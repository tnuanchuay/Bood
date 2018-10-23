package main

import (
	"github.com/nlopes/slack"
	"github.com/xanzy/go-gitlab"
)

type MessageHandler interface{
	Rule(event *slack.MessageEvent) bool
	OnMessageReceive(event *slack.MessageEvent)
}

func CreateHandler(rtm *slack.RTM, botName string, git *gitlab.Client, auth SSHAuth) []MessageHandler{
	return []MessageHandler{
		NewMergeRequestDeployMessageHandler(rtm, botName, git),
		NewContainerKillerMessageHandler(rtm, botName, auth),
	}
}

func HandleMessage(handlers []MessageHandler, event *slack.MessageEvent){
	for _, handler := range handlers{
		if handler.Rule(event){
			handler.OnMessageReceive(event)
			return
		}
	}
}
