package modules

import (
    "strings"
    "gobot/common"
    "bufio"
    "os"
	"math/rand"
    "time"
    "strconv"
)


type Quote struct {
    bot *common.Bot
    quotes []string
    stopped bool
}

func NewQuoteModule(bot *common.Bot) *Quote {
    quote := &Quote {
        bot: bot,
        stopped: false,
    }
    rand.Seed(time.Now().Unix())
    quote.ReadQuotes()
    return quote
}


func (q *Quote) OnMessageRecieved(msg *common.BotMessage) {
    has := strings.HasPrefix(msg.Message, "gobot")
    if has && !q.stopped {
        var str = q.quotes[rand.Intn(len(q.quotes))]
        var answer = msg.User + ": " + str
        q.bot.SendMessage(msg.Channel, answer)
    } else if strings.HasPrefix(msg.Message, "!upd quote") {
        q.ReadQuotes()
    } else if strings.HasPrefix(msg.Message, "!stop quote") { 
        q.stopped = true
    } else if strings.HasPrefix(msg.Message, "!quote ") {
         args := strings.Split(msg.Message, " ")
        
        if len(args) != 2 {
            q.bot.SendMessage(msg.Channel, "Ошибка: Неверное количество аргументов, используйте !quote [номер цитаты]")
            return
        }
        
        num, err := strconv.Atoi(args[1])
        
        if num > len(q.quotes) - 1 {
            q.bot.SendMessage(msg.Channel, "Ошибка: Такой цитаты нет!")
            return
        }
        
        if err != nil {
            q.bot.SendMessage(msg.Channel, "Ошибка: используйте !quote [номер цитаты]")
            return
        }
        
        var answer = msg.User + ": " + q.quotes[num]
        q.bot.SendMessage(msg.Channel, answer)
        
    } else {
        q.stopped = false
    }
}

func (q *Quote) ReadQuotes() {
  file, err := os.Open("quote.txt")
  if err != nil {
    return
  }
  defer file.Close()

  

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  
  q.quotes = lines
}