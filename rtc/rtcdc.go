package rtc

import (
	"io"
	"sync"

	"github.com/pion/webrtc/v4"
)

// ReadWriteCloser wraps *webrtc.DataChannel to satisfy the io.ReadWriteCloser interface.
type ReadWriteCloser struct {
	dc   *webrtc.DataChannel
	pr   *io.PipeReader
	pw   *io.PipeWriter
	once sync.Once
}

// New wraps a *webrtc.DataChannel and returns an io.ReadWriteCloser.
func New(dc *webrtc.DataChannel) io.ReadWriteCloser {
	return Wrap(dc)
}

// Wrap wraps a *webrtc.DataChannel and returns *ReadWriteCloser.
func Wrap(dc *webrtc.DataChannel) *ReadWriteCloser {
	pr, pw := io.Pipe()
	rwc := &ReadWriteCloser{
		dc: dc,
		pr: pr,
		pw: pw,
	}

	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		_, _ = pw.Write(msg.Data)
	})

	dc.OnClose(func() {
		_ = pw.Close()
	})

	return rwc
}

// Read reads received DataChannel bytes.
func (r *ReadWriteCloser) Read(p []byte) (n int, err error) {
	return r.pr.Read(p)
}

// Write sends bytes over the DataChannel.
func (r *ReadWriteCloser) Write(p []byte) (n int, err error) {
	if err := r.dc.Send(p); err != nil {
		return 0, err
	}
	return len(p), nil
}

// Close closes both the pipe and the underlying DataChannel.
func (r *ReadWriteCloser) Close() error {
	var err error
	r.once.Do(func() {
		_ = r.pw.Close()
		_ = r.pr.Close()
		err = r.dc.Close()
	})
	return err
}
