package main

import (
	"log"
	"net"
	"runtime"
	"time"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept rewrite net.Listener interface,
// use TCP Keep-Alive instead
func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

type HandlerInterface interface {
	Handle(conn net.Conn)
}

type Server struct {
	Addr    string
	Handler HandlerInterface
}

func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	var tempDelay time.Duration
	for {
		rw, err := l.Accept()
		if err == nil {
			go srv.ServeConn(rw)
			continue
		}

		// Judge if it is a trmporary error,
		// if it is then sleep for a moment
		if ne, ok := err.(net.Error); ok && ne.Temporary() {
			if tempDelay == 0 {
				tempDelay = 5 * time.Millisecond
			} else {
				tempDelay *= 2
			}
			if max := 1 * time.Second; tempDelay > max {
				tempDelay = max
			}
			log.Printf("Accept error: %v, retrying in %v\n", ne, tempDelay)
			time.Sleep(tempDelay)
			continue
		}
		return err
	}
}

// ListenAndServe start a server for dealing with connection,
// default port is 9102
func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":9102"
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return srv.Serve(tcpKeepAliveListener{l.(*net.TCPListener)})
}

// ServeConn recover handler process
func (srv *Server) ServeConn(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("ServeConn panic serving %v: %v\n%s", conn.RemoteAddr(), err, buf)
			if conn != nil {
				conn.Close()
			}
		}
	}()
	srv.Handler.Handle(conn)
}
