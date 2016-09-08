package common

// BotMessage is an exported type that.
type BotMessage struct {
    Channel    string
    User       string
    Message    string
}

// NewMessage creates and returns objects of
// the exported type Message.
func NewMessage(user, channel, message string) *BotMessage {
    msg := &BotMessage {
        Channel: channel,
        User:    user,
        Message: message,
    }
    return msg
}
