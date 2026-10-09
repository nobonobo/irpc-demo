package main

import (
	"log"

	"irpc-demo/rtc"
)

func main() {
	log.SetFlags(log.Lshortfile)
	rtc.RunHost(rtc.HostID)
}
