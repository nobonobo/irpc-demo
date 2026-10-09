package main

import (
	"log"

	"github.com/nobonobo/irpc-demo/rtc"
	"github.com/nobonobo/irpc-demo/services/rtcmsg"
)

type session struct{}

func (s *session) Send(msg rtcmsg.Message) error {
	log.Println(msg.Name, " -> ", msg.Text)
	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	rtc.RunHost(&session{}, rtc.HostID)
}
