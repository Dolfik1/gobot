package modules

import (
    "strings"
    "gobot/common"
    "net/http"
    "encoding/json"
    "io/ioutil"
)

type OboobsJson struct {
    ID          int `json:"id"`
    Preview     string `json:"preview"`
}

type OboobsModule struct {
    bot *common.Bot
}

func NewOboobsModule(bot *common.Bot) *OboobsModule {
    mod := &OboobsModule {
        bot: bot,
    }
    return mod
}

func (mod *OboobsModule) OnMessageRecieved(msg *common.BotMessage) {
    has := strings.HasPrefix(msg.Message, "!сиськи")
    if has {
        resp, err := http.Get("http://api.oboobs.ru/noise/%7Bcount=1;%20sql%20limit%7D")
        if err == nil {
            defer resp.Body.Close()
            responseData, err := ioutil.ReadAll(resp.Body)
            if err == nil {
                jsonData := &[1]OboobsJson{}
                err = json.Unmarshal(responseData, jsonData)
                
                if err == nil {
                    mod.bot.SendMessage(msg.Channel, "[oboobs] http://media.oboobs.ru/" + jsonData[0].Preview)
                } else {
                    mod.bot.SendMessage(msg.Channel, "Ошибка при парсинге json: " + err.Error())
                }
                
            } else {
                mod.bot.SendMessage(msg.Channel, "Ошибка: " + err.Error())
            }
            
        } else {
            mod.bot.SendMessage(msg.Channel, "Ошибка. Не могу получить сиськи :(") 
        }
    }
}
