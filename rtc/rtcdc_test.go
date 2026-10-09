package rtc_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"irpc-demo/rtc"

	"github.com/pion/webrtc/v4"
)

func TestRTCDCReadWriteCloser(t *testing.T) {
	api := webrtc.NewAPI()
	pc1, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatal(err)
	}
	defer pc1.Close()

	pc2, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatal(err)
	}
	defer pc2.Close()

	pc1.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil {
			_ = pc2.AddICECandidate(c.ToJSON())
		}
	})
	pc2.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil {
			_ = pc1.AddICECandidate(c.ToJSON())
		}
	})

	dc1, err := pc1.CreateDataChannel("test", nil)
	if err != nil {
		t.Fatal(err)
	}

	dc2Chan := make(chan *webrtc.DataChannel, 1)
	pc2.OnDataChannel(func(dc *webrtc.DataChannel) {
		dc2Chan <- dc
	})

	offer, err := pc1.CreateOffer(nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := pc1.SetLocalDescription(offer); err != nil {
		t.Fatal(err)
	}
	if err := pc2.SetRemoteDescription(offer); err != nil {
		t.Fatal(err)
	}

	answer, err := pc2.CreateAnswer(nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := pc2.SetLocalDescription(answer); err != nil {
		t.Fatal(err)
	}
	if err := pc1.SetRemoteDescription(answer); err != nil {
		t.Fatal(err)
	}

	var dc2 *webrtc.DataChannel
	select {
	case dc2 = <-dc2Chan:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for dc2")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan struct{})
	dc1.OnOpen(func() {
		close(done)
	})

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("timeout waiting for dc1 open")
	}

	stream1 := rtc.New(dc1)
	stream2 := rtc.New(dc2)

	defer stream1.Close()
	defer stream2.Close()

	testData := []byte("Hello, WebRTC DataChannel!")

	go func() {
		_, writeErr := stream1.Write(testData)
		if writeErr != nil {
			t.Errorf("write error: %v", writeErr)
		}
	}()

	buf := make([]byte, 1024)
	n, readErr := stream2.Read(buf)
	if readErr != nil {
		t.Fatalf("read error: %v", readErr)
	}

	if !bytes.Equal(buf[:n], testData) {
		t.Fatalf("expected %s, got %s", testData, buf[:n])
	}
}
