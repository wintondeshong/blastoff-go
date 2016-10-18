package lib

import (
	"golang.org/x/net/websocket"
)

type ISlackHandler interface {
	CanHandleMessage(Message) bool
	HandleMessage(*websocket.Conn, Message)
}
