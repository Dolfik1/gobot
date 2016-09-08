package modules

import (
    "encoding/xml"
    "gobot/common"
    "net/http"
    "strings"
    "io/ioutil"
    //"time"
    "fmt"
)

type WeatherModule struct {
    bot              *common.Bot
    WCities          *WeatherCities
    IsCitiesRecieved bool
}

type WeatherCities struct {
    XMLName xml.Name `xml:"cities"`
    Country []*WeatherCountry `xml:"country"`
}

type WeatherCountry struct {
    XMLName     xml.Name `xml:"country"`
    Name        string `xml:"name,attr"`
    City        []*WeatherCity `xml:"city"`
}

type WeatherCity struct {
    XMLName     xml.Name `xml:"city"`
    ID          int `xml:"id,attr"`
    Region      int `xml:"region,attr"`
    Name        string `xml:",chardata"`
}

type WeatherInfo struct {
    XMLName     xml.Name `xml:"info"`
    Weather     *WeatherData `xml:"weather"`
}

type WeatherData struct {
    XMLName     xml.Name `xml:"weather"`
    Day         *WeatherDay `xml:"day"`
}

type WeatherDay struct {
    XMLName     xml.Name `xml:"day"`
    DayParts    []*WeatherDayPart `xml:"day_part"`
}

type WeatherDayPart struct {
    XMLName             xml.Name `xml:"day_part"`
    TypeID              int `xml:"typeid,attr"`
    WeatherType         string `xml:"weather_type"`
    WindSpeed           string `xml:"wind_speed"`
    WindDirection       string `xml:"wind_direction"`
    Dampness            int `xml:"dampness"`
    Pressure            int `xml:"pressure"`
    Temperature         string `xml:"temperature"`    
}

func NewWeatherModule(bot *common.Bot) *WeatherModule {
    module := &WeatherModule {
        bot: bot,
        WCities: &WeatherCities{},
        IsCitiesRecieved: false,
    }
    return module
}

func (w *WeatherModule) OnMessageRecieved(msg *common.BotMessage) {
    has := strings.HasPrefix(msg.Message, "!погода")
    if has {
        args := strings.Split(msg.Message, " ")
        
        if len(args) != 2 {
            w.bot.SendMessage(msg.Channel, "Ошибка: Неверное количество аргументов, используйте !погода [город]")
            return
        }
        cityName := args[1]
        
        if w.IsCitiesRecieved || w.LoadCities() {
            city := w.GetCityByName(cityName)
            if city != nil {
                weather := city.GetWeather()
                if weather != nil {
                    var message = fmt.Sprintf("Погода %s: %s, Температура: %s°C, Ветер %s м/с, %s, Влажность %d%%, Давление %d мм рт. ст.", 
                        city.Name, weather.WeatherType, weather.Temperature, weather.WindSpeed, weather.WindDirection, weather.Dampness, weather.Pressure)
                    w.bot.SendMessage(msg.Channel, message) 
                } else {
                    w.bot.SendMessage(msg.Channel, "Ошибка. Не могу получить погоду.")
                }
            } else {
                w.bot.SendMessage(msg.Channel, "Город \"" + cityName + "\" не найден.")
            }
        }
    }
}

func (w *WeatherModule) LoadCities() bool { 
    resp, err := http.Get("https://pogoda.yandex.ru/static/cities.xml")
    if err == nil {
        defer resp.Body.Close()
        responseData, err := ioutil.ReadAll(resp.Body)
            
        if err == nil {
            err := xml.Unmarshal(responseData, &w.WCities)
                
            if err == nil {
                w.IsCitiesRecieved = true
            } else {
                w.bot.SendMessage(w.bot.Owner, "Ошибка при парсинге xml городов: " + err.Error())
                return false
            }
        } else {
            w.bot.SendMessage(w.bot.Owner, "Ошибка при попытке прочесть ответ сервера: " + err.Error())
            return false
        }
    } else {
        w.bot.SendMessage(w.bot.Owner, "Ошибка при получении городов: " + err.Error())
        return false
    }
        
    return true
}


func (w *WeatherModule) GetCityByName(name string) *WeatherCity { 
    name = strings.ToLower(name)
    if w.IsCitiesRecieved {
        for i := range w.WCities.Country {
            for k := range w.WCities.Country[i].City {
                if strings.ToLower(w.WCities.Country[i].City[k].Name) == name {
                    return w.WCities.Country[i].City[k]
                } 
            }
        }
    }
    return nil
}

func (city *WeatherCity) GetWeather() *WeatherDayPart {
    resp, err := http.Get("http://export.yandex.ru/bar/reginfo.xml?region=" + string(city.Region))
    if err == nil {
        defer resp.Body.Close()
        responseData, err := ioutil.ReadAll(resp.Body)
        if err == nil {
            weather := &WeatherInfo {}
            err := xml.Unmarshal(responseData, &weather)
            if err == nil {
                /*hr := time.Now().Hour()
                
                
                var daypart = 0
                
                if hr <= 6 {
                    daypart = 1
                } else if hr <= 12 {
                    daypart = 2
                } else if hr <= 18 {
                    daypart = 3
                } else {
                    daypart = 4
                }
                
                for i := range weather.Weather.Day.DayParts {
                    if weather.Weather.Day.DayParts[i].TypeID == daypart {
                        dp := weather.Weather.Day.DayParts[i]
                        return dp//weather.Weather.Day.DayParts[i]
                    }
                }
                */
                
                if len(weather.Weather.Day.DayParts) > 0 {
                    return weather.Weather.Day.DayParts[0]
                }
                return nil
            }
            fmt.Println("Ошибка. Не могу распарсить xml погоды: " + err.Error())        
        } else {
           fmt.Println("Ошибка. Не могу прочитать ответ погоды от сервера: " + err.Error());
        }
    } else {
        fmt.Println("Возникла ошибка при запросе погоды: " + err.Error());
    }
    return nil
}
