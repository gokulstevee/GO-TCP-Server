package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	address string
	msg     []byte
}

type Server struct {
	listnerAddr string
	ln          net.Listener
	quitch      chan struct{}
	msgchan     chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listnerAddr: listenAddr,
		quitch:      make(chan struct{}),
		msgchan:     make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listnerAddr)
	if err != nil {
		return err
	}

	defer ln.Close()

	s.ln = ln

	s.acceptLoop()

	<-s.quitch

	close(s.msgchan)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		fmt.Println("conn...", conn)
		if err != nil {
			fmt.Println("Accept Error: ", err)
			continue
		}

		fmt.Println("New connection to the server: ", conn.RemoteAddr())

		go s.readLoop(conn)

	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 2048)

	fmt.Println("n...", conn)

	for {
		n, err := conn.Read(buff)

		fmt.Println("n...", n)

		if err != nil {
			fmt.Println("Read buff Error: ", err)
			continue
		}

		s.msgchan <- Message{
			address: conn.RemoteAddr().String(),
			msg:     buff[:n],
		}

		conn.Write([]byte("Thank you!\n"))

	}
}

func main() {
	server := NewServer(":3000")

	go func() {
		for msg := range server.msgchan {
			fmt.Println("A new message received, ", msg.address, string(msg.msg))
		}
	}()

	log.Fatal(server.Start())
}
