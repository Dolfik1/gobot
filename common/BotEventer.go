package common

type BotEventer interface {
	OnMessageRecieved(message *BotMessage)
}
