package main

import (
	"context"
	"log"
	"os"
	"time"
	"uuid"

	"irpc-demo/rtc"
	"irpc-demo/services/rtcmsg"
)

func main() {
	name := "Unknown"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	log.SetFlags(log.Lshortfile)
	id := uuid.NewV4()
	ctx := context.Background()
	ch := make(chan rtcmsg.Message)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			ch <- rtcmsg.Message{
				Name: name,
				Text: "hello",
			}
		}
	}()
	rtc.RunNode(ctx, rtc.HostID, id.String(), name, ch)
}
