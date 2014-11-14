package main

import (
	"io"
	"log"
	"net"
	"runtime"

	"thriftAndGob/trans"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(conn net.Conn) {
	cmr, err := trans.NewConsumer(conn, false)
	if err != nil {
		log.Printf("NewConsumer error: %v\n", err)
		return
	}
	defer cmr.Close()
	for {
		msgs, err := cmr.Receive()
		if err != nil {
			if err == io.EOF {
				log.Printf("Close Conn")
				return
			}
			log.Printf("Receive error: %v\n", err)
			return
		}
		_ = msgs
	}
}

func main() {
	server := &Server{":9102", NewHandler()}
	log.Println(server.ListenAndServe())
}
