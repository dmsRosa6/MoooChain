package network

import (
	"fmt"
	"time"
)

type ServerOpts struct {
	Transport []Transport
}

type Server struct {
	ServerOpts
	rpcChan  chan RPC
	quitChan chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		rpcChan:    make(chan RPC),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(5 * time.Second)

	running := true
	for running {
		select {

		case rpc := <-s.rpcChan:
			fmt.Printf("HANDLE: %+v", rpc)

		case <-s.quitChan:
			running = false

		case <-ticker.C:
			fmt.Println("Tick!!")

		default:

		}
	}

	fmt.Println("Server Started")
}

func (s *Server) initTransports() {
	for _, tr := range s.Transport {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcChan <- rpc
			}
		}(tr)
	}
}
