package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/nobonobo/irpc-demo/services/chat"

	"github.com/marben/irpc"
)

type worker struct {
}

func (w *worker) OnMessage(msg chat.Message) {
	fmt.Printf("\r%s: %q\n", msg.Name, msg.Text)
	fmt.Print("> ")
}

func main() {
	name := "unknown"
	flag.StringVar(&name, "name", name, "name of worker")
	flag.Parse()
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	workerImpl := &worker{}
	event := chat.NewEventIrpcService(workerImpl)
	ep := irpc.NewEndpoint(conn, irpc.WithEndpointServices(event))
	defer ep.Close()

	client, err := chat.NewSessionIrpcClient(ep)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		line := scanner.Text()
		if err := client.Send(chat.Message{
			Name: name,
			Text: line,
		}); err != nil {
			log.Fatal(err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
