package akevitt

import "github.com/globalcyberalliance/telnet-go"

type telnetServer struct {
	addr string
}

func (t *telnetServer) Run(e *Engine) error {
	defer e.wg.Done()
	return telnet.ListenAndServe(t.addr, func(server *telnet.Session) {
		ctx := &Context{
			addr: server.RemoteAddr(),
			rw:   server,
		}

		e.AddSession(ctx)

		e.InvokeContext(ctx)
	})
}

func (t *telnetServer) Name() string {
	return "Telnet"
}

func Telnet(addr string) Handler {
	return &telnetServer{addr: addr}
}
