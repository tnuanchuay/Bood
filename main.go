package main

import (
	"github.com/nlopes/slack"
	"fmt"
	"log"
	"os"
	"github.com/xanzy/go-gitlab"
	"io/ioutil"
	"encoding/json"
)

func main(){
	configuration := ReadConfiguration()
	sshAuth := SSHAuth{configuration.SSHUsername, configuration.SSHPassword}

	api := slack.New(
		configuration.SlackApiKey,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	rtm := api.NewRTM()

	git := gitlab.NewClient(nil, configuration.GitlabApiKey)
	git.SetBaseURL("http://gitlab.grafana48.com/api/v4")

	pipelineInstances := CreatePipelineInstance()
	pipeline := CreatePipeline(pipelineInstances)
	handlers := CreateHandler(rtm, configuration.BotName, git, sshAuth)

	go rtm.ManageConnection()
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace C2147483705 with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C2147483705"))

		case *slack.MessageEvent:
			pipelinedEvent := pipeline(ev)
			if pipelinedEvent != nil{
				go HandleMessage(handlers, pipelinedEvent)
			}

		default:
		}
	}
}

func ReadConfiguration() Configuration{
	b, err := ioutil.ReadFile("./configuration.json")
	if err != nil {
		panic(err)
	}

	config := Configuration{}
	err = json.Unmarshal(b, &config)
	if err != nil{
		panic(err)
	}
	return config
}

type SSHAuth struct{
	Username string
	Password string
}

type Configuration struct{
	SlackApiKey		string
	GitlabApiKey	string
	SSHUsername		string
	SSHPassword		string
	BotName			string
}
