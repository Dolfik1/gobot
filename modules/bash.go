package modules

import (
    "gobot/common"
    "strings"
    "net/http"
    "io/ioutil"
    "golang.org/x/net/html"
    "bytes"
    "gopkg.in/xmlpath.v2"
    "golang.org/x/text/encoding/charmap"
)

type BashModule struct {
    bot *common.Bot
}

func NewBashModule(bot *common.Bot) *BashModule {
    mod := &BashModule {
        bot: bot,
    }
    return mod
}


func (mod *BashModule) OnMessageRecieved(msg *common.BotMessage) { 
    has := strings.HasPrefix(msg.Message, "!bash") || strings.HasPrefix(msg.Message, "!bashim") || strings.HasPrefix(msg.Message, "!баш")
    if has {
        resp, err := http.Get("http://bash.im/random")
        if err == nil {
            defer resp.Body.Close()
            responseData, err := ioutil.ReadAll(resp.Body)
            if err == nil {
                
                decoder := charmap.Windows1251.NewDecoder()
                responseData, err = decoder.Bytes(responseData)
                
                reader := strings.NewReader(string(responseData))
                root, err := html.Parse(reader)

                if err == nil {
                    
                    var b bytes.Buffer
                    html.Render(&b, root)
                    fixedHTML := b.String()

                    reader = strings.NewReader(fixedHTML)
                    xmlroot, err := xmlpath.ParseHTML(reader)

                    if err == nil {                        
                        path := xmlpath.MustCompile("//body/div[2]/div[3]/div[2]")
                        if value, ok := path.String(xmlroot); ok {
                            
                            strings.Replace(value, "<br></br>", "\n", -1)
                            strings.Replace(value, "&quot;", "\"", -1)
                            strings.Replace(value, "<br>", "\n", -1)
                            strings.Replace(value, "<br />", "\n", -1)
                            strings.Replace(value, "&lt;", "<", -1)
                            strings.Replace(value, "&gt;", ">", -1)
                            
                            mod.bot.SendMessage(msg.Channel, value)
                        }
                    } else {
                        mod.bot.SendMessage(msg.Channel, "Ошибка при парсинге HTML: " + err.Error())
                    }
                } else {
                    mod.bot.SendMessage(msg.Channel, "Ошибка при парсинге HTML: " + err.Error())
                }
            } else {
                mod.bot.SendMessage(msg.Channel, "Ошибка: " + err.Error())
            }
        } else {
            mod.bot.SendMessage(msg.Channel, "Ошибка. Не могу получить сиськи :(") 
        }
    }
}