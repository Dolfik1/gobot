package modules

import (
    "strings"
    "gobot/common"
)

// Fake module
type Fake struct {
    bot *common.Bot
}

// NewFakeModule create fake module
func NewFakeModule(bot *common.Bot) *Fake {
    fake := &Fake {
        bot: bot,
    }
    return fake
}

// OnMessageRecieved called, when bot receive message
func (f *Fake) OnMessageRecieved(msg *common.BotMessage) {
    //gobot_
    //!fake [channel] [user] [message]
    has := strings.HasPrefix(msg.Message, "!fake ")
    if has {
        args := strings.Split(msg.Message, " ")
        
        if len(args) < 3 {
            //todo error
            return
        }
        start := len(args[0]) + len(args[1]) + len(args[2]) + 3// 3 - пробелы
        text := msg.Message[start:len(msg.Message)]
        m := common.NewMessage(args[2], args[1], text)
        f.bot.MessageRecieved(m)
    }
}