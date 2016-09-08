package common

import (
	"fmt"
	"github.com/thoj/go-ircevent"
)

type Bot struct {
	serverAddress string
	botNickname   string
	channels      []string
	irc           *irc.Connection
	modules	 	  []BotEventer
	Owner		  string
}

func BOT(server, nick, owner string, channels []string) *Bot {
	if len(server) == 0 {
		return nil
	}

	if len(nick) == 0 {
		return nil
	}

	bot := &Bot{
		botNickname:   nick,
		serverAddress: server,
		channels:      channels,
		modules:	   []BotEventer { },
		Owner: 		   owner,
	}


	bot.irc = irc.IRC(bot.botNickname, bot.botNickname)

	err := bot.irc.Connect(bot.serverAddress)
	if err != nil {
		fmt.Println("Failed connecting")
	}

	bot.irc.AddCallback("001", func(e *irc.Event) {
		for idx := range bot.channels {
			bot.irc.Join(bot.channels[idx])
			fmt.Println("Joined to the channel: " + bot.channels[idx])
		}
	})

	bot.irc.AddCallback("PRIVMSG", func(event *irc.Event) {
		//event.Message() contains the message
		//event.Nick Contains the sender
		//event.Arguments[0] Contains the channel
		msg := NewMessage(event.Nick, event.Arguments[0], event.Message())
		for i := range bot.modules {
			go bot.modules[i].OnMessageRecieved(msg)
		}
	});
	

	return bot
}

func (bot *Bot) SendMessage(messageTo, message string) { 
	bot.irc.Privmsg(messageTo, message)
}

func (bot *Bot) RegisterModule(module BotEventer) {
	bot.modules = append(bot.modules, module)
}

func (bot *Bot) MessageRecieved(msg *BotMessage) {
	for i := range bot.modules {
		go bot.modules[i].OnMessageRecieved(msg)
	}
}

func (bot *Bot) Loop() {
	bot.irc.Loop()
}