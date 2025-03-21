package akevitt

import (
	"io"
	"net"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"
)

type Context struct {
	addr net.Addr
	rw   io.ReadWriter
	tty  tty

	width  int
	height int
}

func (ctx *Context) Write(p []byte) (n int, err error) {
	return ctx.rw.Write(p)
}

func (ctx *Context) Read(p []byte) (n int, err error) {
	return ctx.rw.Read(p)
}

func (ctx *Context) ClientIP() net.Addr {
	return ctx.addr
}

type tty struct {
	size struct {
		Width  int
		Height int
	}
	ch       <-chan ssh.Window
	resizecb func()
	mu       sync.Mutex
}

func (t *tty) Start() error {
	go func() {
		for win := range t.ch {
			t.size = win
			t.notifyResize()
		}
	}()
	return nil
}
func (t *tty) Stop() error {
	return nil
}
func (t *tty) Drain() error {
	return nil
}
func (t *tty) WindowSize() (tcell.WindowSize, error) {
	return tcell.WindowSize{
		Width:       t.size.Width,
		Height:      t.size.Height,
		PixelWidth:  16,
		PixelHeight: 16,
	}, nil
}
func (t *tty) NotifyResize(cb func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.resizecb = cb
}
func (t *tty) notifyResize() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.resizecb != nil {
		t.resizecb()
	}
}

func (t *tty) Close() error {
	return nil
}
