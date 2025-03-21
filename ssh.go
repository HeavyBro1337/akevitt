package akevitt

import "github.com/gliderlabs/ssh"

type sshServer struct {
	addr    string
	rsaFile string
}

func (s *sshServer) Name() string {
	return "SSH"
}

func (s *sshServer) Run(e *Engine) error {
	usePubKey := ssh.HostKeyFile(s.rsaFile)

	allowKeys := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		return true
	})

	return ssh.ListenAndServe(s.addr, func(s ssh.Session) {
		ctx := &Context{
			addr: s.RemoteAddr(),
			rw:   s,
		}

		e.AddSession(ctx)

		e.InvokeContext(ctx)
	}, allowKeys, usePubKey)
}

func SSH(addr, keyPath string) Handler {
	return &sshServer{
		addr:    addr,
		rsaFile: keyPath,
	}
}
