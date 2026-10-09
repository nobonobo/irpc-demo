package main

import (
	"log"

	"github.com/nobonobo/irpc-demo/rtc"
)

func main() {
	log.SetFlags(log.Lshortfile)
	rtc.RunHost(rtc.HostID)
}
