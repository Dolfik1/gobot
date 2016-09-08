package modules

import (
    "strings"
    "gobot/common"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "fmt"
    "strconv"
    "time"
)

// Cryptocurrency bot module
type Cryptocurrency struct {
    bot         *common.Bot
    triggers    []CryptocurrencyTriggerModel
}

// CryptocurrencyModelTicker is model from json server answer
type CryptocurrencyModelTicker struct {
    Base        string      `json:"base"`
    Target      string      `json:"target"`
    Price       string      `json:"price"`
    Volume      string      `json:"volume"`
    Change      string      `json:"change"`
    Error       string      `json:"error"`
    Success     bool        `json:"success"`
}

// CryptocurrencyModel is model from json server answer
type CryptocurrencyModel struct {
    Error       string      `json:"error"`
    Success     bool        `json:"success"`
    Timestamp   int         `json:"timestamp"`
    Ticker*     CryptocurrencyModelTicker `json:"ticker"`
}

// CryptocurrencyTriggerModel - model
type CryptocurrencyTriggerModel struct {
    Base                string
    Target              string
    Value               float64
    RequestedChannel    string
}

// NewCryptocurrency create new cryptocurrency module object
func NewCryptocurrency(bot *common.Bot) *Cryptocurrency {
    cryptocurrency := &Cryptocurrency {
        bot: bot,
        triggers: []CryptocurrencyTriggerModel {},
    }
    go cryptocurrency.UpdateTriggers()
    return cryptocurrency
}

// OnMessageRecieved called, when bot receive message
func (module *Cryptocurrency) OnMessageRecieved(msg *common.BotMessage) {
    has := strings.HasPrefix(msg.Message, "!rate ")
    if has {
        module.RateCommand(msg)
    }

    if strings.HasPrefix(msg.Message, "!trigger ") {
        module.TriggerCommand(msg)
    }

    if strings.HasPrefix(msg.Message, "!triggers") {
        module.TriggersListCommand(msg)
    }
}

// TriggersListCommand - process !triggres command
func (module *Cryptocurrency) TriggersListCommand(msg *common.BotMessage) {

    if len(module.triggers) == 0 {
        module.bot.SendMessage(msg.Channel, "Список триггеров пустой")
        return
    }

    result := ""

    for i := range module.triggers {
        trig := &module.triggers[i]
        result += fmt.Sprintf("%s/%s - %f, ", trig.Base, trig.Target, trig.Value)
	}

    module.bot.SendMessage(msg.Channel, result)
} 

// TriggerCommand - !trigger [base] [target] [change]
func (module *Cryptocurrency) TriggerCommand(msg *common.BotMessage) {
    args := strings.Split(msg.Message, " ")
    if len(args) != 4 {
        module.bot.SendMessage(msg.Channel, "Ошибка: Неверное количество аргументов, используйте !trigger [base] [target] [change]")
        return
    }

    change, err := strconv.ParseFloat(args[3], 32)
    if err != nil {
        module.bot.SendMessage(msg.Channel, "Ошибка: Неверное значение аргумента change, значением должно быть число. Используйте !trigger [base] [target] [change]")
        return
    }

    trigger := CryptocurrencyTriggerModel{}
    trigger.Base = args[1]
    trigger.Target = args[2]
    trigger.Value = change
    trigger.RequestedChannel = msg.Channel

    isTriggerExists := false

    for i := range module.triggers {
		trig := &module.triggers[i]
        if trig.Base == trigger.Base && trig.Target == trigger.Target {
            trig.Base = trigger.Base
            trig.Target = trigger.Target
            isTriggerExists = true
            break
        }
	}

    if isTriggerExists {
       return 
    }

    module.triggers = append(module.triggers, trigger)   
}

// RateCommand outputs rate
func (module *Cryptocurrency) RateCommand(msg *common.BotMessage) {
    args := strings.Split(msg.Message, " ")
        
    if len(args) != 3 && len(args) != 4 {
        module.bot.SendMessage(msg.Channel, "Ошибка: Неверное количество аргументов, используйте !rate [количество = 1] [валюта 1] [валюта 2]")
        return
    }

    curr1 := ""
    curr2 := ""
    count := 1.0

   if len(args) == 3 {
       curr1 = args[1]
       curr2 = args[2]
    } else {
        curr1 = args[2]
        curr2 = args[3]
        cnt, err := strconv.ParseFloat(args[1], 32)
        if err != nil {
            module.bot.SendMessage(msg.Channel, "Ошибка: Неверный формат первого аргумента, используйте !rate [количество = 1] [валюта 1] [валюта 2]")
            return
        }
        count = cnt
    }

    model, err := module.GetTicker(curr1, curr2)

    if err != nil {
        return
    }

    if model.Success == true {
        priceFloat, err := strconv.ParseFloat(model.Price, 64)
        if err != nil {
            module.bot.SendMessage(msg.Channel, "Ошибка при парсинге цены: " + model.Price + ", " + err.Error())
            return
        } 
        sum := count * priceFloat
        msgText := fmt.Sprintf("%f %s = %f %s, %s", count, model.Base, sum, model.Target, model.Change)
        module.bot.SendMessage(msg.Channel, msgText)//count + " " + model.Ticker.Base + " = " + model.Ticker.Price + " " + model.Ticker.Target + ", " + model.Ticker.Change)
    } else {
        module.bot.SendMessage(msg.Channel, "Ошибка: " + model.Error)
    }
}

// GetTicker returns ticker from api
func (module *Cryptocurrency) GetTicker(curr1, curr2 string) (CryptocurrencyModelTicker, error) {
    model := &CryptocurrencyModel { }

    url := "https://www.cryptonator.com/api/full/" + curr1 + "-" + curr2;
    resp, err := http.Get(url)

    if err != nil {
        module.bot.SendMessage(module.bot.Owner, "Ошибка при запросе валюты: " + err.Error())
        fmt.Print("Ошибка при запросе валюты: " + err.Error())
        return *model.Ticker, err
    }

    defer resp.Body.Close()
    responseData, err := ioutil.ReadAll(resp.Body)
    text := string(responseData)
    fmt.Print(text)

    if err != nil {
        module.bot.SendMessage(module.bot.Owner, "Ошибка при попытке прочесть ответ сервера: " + err.Error())
        fmt.Print("Ошибка при попытке прочесть ответ сервера: " + err.Error())
        return *model.Ticker, err
    }

    err = json.Unmarshal(responseData, &model)

    if err != nil {
        module.bot.SendMessage(module.bot.Owner, "Ошибка при парсинге JSON ответа: " + err.Error())
        fmt.Print("Ошибка при парсинге JSON ответа: " + err.Error())
        return *model.Ticker, err
    }

    model.Ticker.Success = true
    return *model.Ticker, nil
}

// UpdateTriggers - check updates in triggers
func (module *Cryptocurrency) UpdateTriggers() {
    for i := range module.triggers {
		trig := &module.triggers[i]
        curr, err := module.GetTicker(trig.Base, trig.Target)
        if err != nil {
            continue
        }
        
        change, err := strconv.ParseFloat(curr.Change, 32)
        if err != nil {
            continue
        }

        if change >= trig.Value {
            module.bot.SendMessage(trig.RequestedChannel, "Цена " + trig.Base + "/" + trig.Target + " изменилось. Новое значение: " + curr.Change)
        }

        time.Sleep(10 * time.Second)
	}

    time.Sleep(5 * time.Minute)
    go module.UpdateTriggers()
}