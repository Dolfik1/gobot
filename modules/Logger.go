package modules

import (
    "strings"
    "gobot/common"
    "fmt"
    "time"
)

type Logger struct {
    bot *common.Bot
}

func LOGGER(bot *common.Bot) *Logger {
    logger := &Logger {
        bot: bot,
    }
    return logger
}

func (l *Logger) OnMessageRecieved(msg *common.BotMessage) {
    has := strings.HasPrefix(msg.Message, "!привет")
    if has {
        l.bot.SendMessage(msg.Channel, "Привет " + msg.User + "!")
    }
    fmt.Printf("[%s] %s: %s\n", time.Now(), msg.User, msg.Message)
}