package messagehandler

type IMessageHandler interface {
	HandleMessage([]byte)
}
