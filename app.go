package main

import (
	"fmt"
	"log"
	"net"
)

// 127.0.0.1:7878

// type Server interface {
// 	start()
// }

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	listner    net.Listener
	channel    chan struct{}
	msgChannel chan Message
}

func NewServer(listendAddr string) *Server {
	return &Server{
		listenAddr: listendAddr,
		channel:    make(chan struct{}),
		msgChannel: make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.listner = ln

	s.acceptConnection()

	fmt.Println("Server listening at port: ", s.listenAddr)
	<-s.channel
	close(s.msgChannel)
	return nil
}

func (s *Server) acceptConnection() {
	for {
		conn, err := s.listner.Accept()
		if err != nil {
			log.Fatal("Connection error: ", err)
			continue
		}

		fmt.Println("Connection established with:", conn.RemoteAddr())
		conn.Write([]byte("Connection established, welcome to gtc_server\n"))

		go s.processConnection(conn)

	}
}

func (s *Server) processConnection(conn net.Conn) {
	defer conn.Close()

	buff := make([]byte, 2048)

	for {
		n, err := conn.Read(buff)
		if err != nil {
			log.Fatal("Connection Error: ", err)
			continue
		}

		s.msgChannel <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buff[:n],
		}
		// msg := buff[:n]
		// fmt.Println(string(msg))
	}
}

func main() {
	server := NewServer(":8000")

	go func() {
		for msg := range server.msgChannel {
			fmt.Printf("Received message from connection %s: %s", msg.from, string(msg.payload))
		}
	}()

	log.Fatal(server.Start())

}
