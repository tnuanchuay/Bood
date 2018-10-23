package main

import "github.com/nlopes/slack"

type EchoMessagePipeline struct{
	Botname	string
}

func NewEchoMessagePipeline(botname string) EchoMessagePipeline {
	return EchoMessagePipeline{Botname:botname}
}

func (e EchoMessagePipeline) Do(ev *slack.MessageEvent) *slack.MessageEvent{
	if ev == nil{
		return nil
	}

	if ev.Username == e.Botname {
		return nil
	}

	return ev
}