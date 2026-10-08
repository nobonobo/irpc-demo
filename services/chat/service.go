package chat

//go:generate irpc $GOFILE

type Message struct {
	Name string
	Text string
}

type Event interface {
	OnMessage(msg Message)
}

type Session interface {
	Send(msg Message) error
}
