package rtc

import (
	"context"
	"log"

	"github.com/marben/irpc"
	"github.com/nobonobo/rtcconnect/node"
	"github.com/pion/webrtc/v4"

	"github.com/nobonobo/irpc-demo/services/rtcmsg"
)

func RunHost(impl rtcmsg.Session, id string) {
	host := node.NewHost(id)
	defer host.Close()
	host.OnConnected = func(n *node.Node) {
		log.Printf("Connected: %s", n.ID())
		n.PeerConnection().OnDataChannel(func(dc *webrtc.DataChannel) {
			conn := New(dc)
			svc := rtcmsg.NewSessionIrpcService(impl)
			ep := irpc.NewEndpoint(conn, irpc.WithEndpointServices(svc))
			log.Println("ep:", ep)
		})
	}
	defer host.Close()
	ctx := context.Background()
	log.Println("host started:", id)
	if err := host.Listen(ctx); err != nil {
		log.Fatal(err)
	}
}

func RunNode(ctx context.Context, host, id, name string, ch <-chan rtcmsg.Message) {
	n := node.New(id)
	defer n.Close()
	if err := n.Connect(ctx, host); err != nil {
		log.Fatal(err)
	}
	n.DataChannel().OnOpen(func() {
		log.Println("data channel opened:", host)
		conn := New(n.DataChannel())
		ep := irpc.NewEndpoint(conn)
		client, _ := rtcmsg.NewSessionIrpcClient(ep)
		log.Println("client:", client)
		go func() {
			for msg := range ch {
				if err := client.Send(msg); err != nil {
					log.Println("error:", err)
				}
			}
		}()
	})
	n.DataChannel().OnClose(func() {
		log.Println("data channel closed:", host)
	})
	n.DataChannel().OnMessage(func(msg webrtc.DataChannelMessage) {
		log.Println("received:", string(msg.Data))
	})
	log.Println("node started:", id)
	<-ctx.Done()
}
