package main

import (
	"log"
	"net"
	"sync"

	"irpc-demo/services/chat"

	"github.com/marben/irpc"
)

type worker struct {
	event chat.Event
}

var (
	mutex   sync.RWMutex
	clients = map[*worker]struct{}{}
)

func (w *worker) Send(msg chat.Message) error {
	mutex.RLock()
	defer mutex.RUnlock()
	for c := range clients {
		c.event.OnMessage(msg)
	}
	return nil
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	// 接続ごとに生成された Endpoint を取得する
	server := irpc.NewServer(
		irpc.WithOnConnect(func(ep *irpc.Endpoint) {
			// このEndpointを通じてClient側のOnMessageを呼ぶ
			event, err := chat.NewEventIrpcClient(ep)
			if err != nil {
				log.Println(err)
				return
			}

			workerImpl := &worker{
				event: event,
			}
			mutex.Lock()
			defer mutex.Unlock()
			clients[workerImpl] = struct{}{}

			// Server側のWorkerサービスを登録
			ep.RegisterService(
				chat.NewSessionIrpcService(workerImpl),
			)
			go func() {
				<-ep.Context().Done()
				mutex.Lock()
				defer mutex.Unlock()
				delete(clients, workerImpl)
			}()
		}),
	)

	log.Fatal(server.Serve(ln))
}
