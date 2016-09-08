package modules

import (
    "strings"
    "gobot/common"
    "github.com/robertkrimen/otto"
)

// JsEval module
type JsEval struct {
    bot *common.Bot
    vm *otto.Otto
}

// NewJsEvalModule create jseval module
func NewJsEvalModule(bot *common.Bot) *JsEval {
    jsEval := &JsEval {
        bot: bot,
        vm: otto.New(),
    }
    return jsEval
}

// OnMessageRecieved called, when bot receive message
func (module *JsEval) OnMessageRecieved(msg *common.BotMessage) {
    //gobot_
    //!fake [channel] [user] [message]
    has := strings.HasPrefix(msg.Message, "!js ")
    if has {
        args := strings.Split(msg.Message, " ")
        
        if len(args) < 2 {
            module.bot.SendMessage(msg.Channel, "Ошибка: Неверное количество аргументов, используйте !js [скрипт]")
            return
        }

        start := len(args[0]) + 1 // 1 - пробел
        code := msg.Message[start:len(msg.Message)]
        val, err := module.vm.Run(code)
        if err != nil {
            module.bot.SendMessage(msg.Channel, "Ошибка: " + err.Error())
            return
        }
        
        str, err := val.ToString()
        if err != nil {
            module.bot.SendMessage(msg.Channel, "Ошибка: Не могу преобразовать в строку: " + err.Error())
            return
        }

        module.bot.SendMessage(msg.Channel, str)
    }
}