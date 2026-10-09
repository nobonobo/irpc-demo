package rtcmsg

//go:generate irpc $GOFILE

type Message struct {
	Name string
	Text string
}

type Session interface {
	Send(msg Message) error
}
