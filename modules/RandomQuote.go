package modules

import (
    "gobot/common"
	"math/rand"
    "time"
    "strings"
)


type RandomQuote struct {
    bot *common.Bot
}

func NewRandomQuoteModule(bot *common.Bot) *RandomQuote {
    quote := &RandomQuote {
        bot: bot,
    }
    rand.Seed(time.Now().Unix())
    return quote
}

func (q *RandomQuote) OnMessageRecieved(msg *common.BotMessage) {
    if strings.Index(msg.Message, "gobot") == -1 && strings.Index(msg.User, "gobot") == -1  {
        n := randInt(0, 1)
        if n == 0 {
            msg.Message = "gobot"
            q.bot.MessageRecieved(msg)
        }
    }
}


func randInt(min, max int) int {
    return min + rand.Intn(max-min)
}