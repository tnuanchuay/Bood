package main

import (
	"github.com/nlopes/slack"
	"os"
	"strings"
	"strconv"
	"golang.org/x/crypto/ssh"
	"fmt"
	"io"
	"net"
)

type ContainerKillerMessageHandler struct{
	botName string
	rtm     *slack.RTM
	username string
	password string
}

func NewContainerKillerMessageHandler(rtm *slack.RTM, botName string, auth SSHAuth) ContainerKillerMessageHandler{
	return ContainerKillerMessageHandler{
		rtm:     rtm,
		botName: botName,
		username:auth.Username,
		password:auth.Password,
	}
}

func (m ContainerKillerMessageHandler) Rule(event *slack.MessageEvent) bool{
	message := event.Text
	splitedMessage := strings.Split(message, " ")

	if len(splitedMessage) != 3{
		return false
	}

	botName, kill, port := func()(string, string, string){
		return splitedMessage[0], splitedMessage[1], splitedMessage[2]
	}()

	if botName != m.botName{
		return false
	}

	if kill != "kill"{
		return false
	}

	_, err := strconv.Atoi(port)
	if err != nil{
		return false
	}

	return true
}

func (m ContainerKillerMessageHandler) OnMessageReceive(event *slack.MessageEvent){
	_, _, port := func()(string, string, string){
		message := event.Text
		splitedMessage := strings.Split(message, " ")
		return splitedMessage[0], splitedMessage[1], splitedMessage[2]
	}()

	sshConfig := &ssh.ClientConfig{
		User: m.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(m.password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	connection, err := ssh.Dial("tcp", "167.99.79.65:22", sshConfig)
	if err != nil {
		m.rtm.SendMessage(m.rtm.NewOutgoingMessage(err.Error(), event.Channel))
		return
	}
	defer connection.Close()

	session, err := connection.NewSession()
	if err != nil {
		m.rtm.SendMessage(m.rtm.NewOutgoingMessage(err.Error(), event.Channel))
		return
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		m.rtm.SendMessage(m.rtm.NewOutgoingMessage(err.Error(), event.Channel))
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		m.rtm.SendMessage(m.rtm.NewOutgoingMessage(err.Error(), event.Channel))
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		m.rtm.SendMessage(m.rtm.NewOutgoingMessage(err.Error(), event.Channel))
	}
	go io.Copy(os.Stderr, stderr)
	command := fmt.Sprintf("/home/sg-ciuser/script.sh %s", port)
	err = session.Run(command)
	if err != nil {
		m.rtm.SendMessage(m.rtm.NewOutgoingMessage(err.Error(), event.Channel))
		return
	}

	m.rtm.SendMessage(m.rtm.NewOutgoingMessage("killing container has been done by command " + command, event.Channel))
}