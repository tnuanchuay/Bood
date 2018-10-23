package main

import (
	"github.com/nlopes/slack"
	"strings"
	"strconv"
	"github.com/xanzy/go-gitlab"
	"fmt"
)

type MergeRequestDeployMessageHandler struct{
	botName string
	rtm     *slack.RTM
	git *gitlab.Client
}

func NewMergeRequestDeployMessageHandler(rtm *slack.RTM, botName string, git *gitlab.Client) MergeRequestDeployMessageHandler{
	return MergeRequestDeployMessageHandler{
		rtm:     rtm,
		botName: botName,
		git:git,
	}
}

func (m MergeRequestDeployMessageHandler) Rule(event *slack.MessageEvent) bool{
	message := event.Text
	splitedMessage := strings.Split(message, " ")

	if len(splitedMessage) != 4{
		return false
	}

	botName, deploy, mr, port := func()(string, string, string, string){
		return splitedMessage[0], splitedMessage[1], splitedMessage[2], splitedMessage[3]
	}()

	if botName != m.botName{
		return false
	}

	if deploy != "deploy"{
		return false
	}

	_, err := strconv.Atoi(mr)
	if err != nil{
		return false
	}

	_, err = strconv.Atoi(port)
	if err != nil{
		return false
	}

	return true
}

func (m MergeRequestDeployMessageHandler) OnMessageReceive(event *slack.MessageEvent){

	mrNumber, port := func()(int, int){
		splitedMessage := strings.Split(event.Text, " ")
		number, _ := strconv.Atoi(splitedMessage[2])
		port, _ := strconv.Atoi(splitedMessage[3])
		return number, port
	}()

	mr, _, err := m.git.MergeRequests.GetMergeRequest(3, mrNumber)
	if err != nil {
		m.rtm.SendMessage(m.rtm.NewOutgoingMessage(err.Error(), event.Channel))
		return
	}

	p, _, err := m.git.Pipelines.CreatePipeline(
		3,
		&gitlab.CreatePipelineOptions{
			Ref:       &mr.SourceBranch,
			Variables: []*gitlab.PipelineVariable{
				&gitlab.PipelineVariable{Key:"BOT", Value:m.botName},
				&gitlab.PipelineVariable{Key:"PORT", Value:fmt.Sprintf("%d", port)},
			},
			})
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println(p)
	}

	m.rtm.SendMessage(m.rtm.NewOutgoingMessage(fmt.Sprintf("deploying merge request %d ที่ port %d", mrNumber, port), event.Channel))
}