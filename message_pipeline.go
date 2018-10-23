package main

import "github.com/nlopes/slack"

type MessagePipeline interface{
	Do(ev *slack.MessageEvent) *slack.MessageEvent
}

func CreatePipeline(pipelines []MessagePipeline) func (event *slack.MessageEvent)*slack.MessageEvent{
	return func (event *slack.MessageEvent)*slack.MessageEvent{
		var recursionEvent *slack.MessageEvent
		recursionEvent = event
		for _, mp := range pipelines{
			recursionEvent = mp.Do(recursionEvent)
		}
		return recursionEvent
	}
}

func CreatePipelineInstance(botname string) []MessagePipeline{
	return []MessagePipeline{
		NewEchoMessagePipeline(botname),
	}
}
