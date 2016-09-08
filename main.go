package main

import (
	"gobot/modules"
	"gobot/common"
)

func main() {
	bot := common.BOT("irc.quakenet.org:6667", "gobot", "#Dolfik", []string{ "#dolfik", "#Lover.ee" })
	bot.RegisterModule(modules.LOGGER(bot))
	bot.RegisterModule(modules.NewWeatherModule(bot))
	bot.RegisterModule(modules.NewQuoteModule(bot))
	bot.RegisterModule(modules.NewFakeModule(bot))
	bot.RegisterModule(modules.NewCryptocurrency(bot))
	bot.RegisterModule(modules.NewJsEvalModule(bot))
	//bot.RegisterModule(modules.NewRandomQuoteModule(bot))
	bot.Loop();
}
