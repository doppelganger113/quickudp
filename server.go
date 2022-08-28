package quickudp

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

type Handler func(msg Message, w Writer)

type Writer interface {
	WriteToUDP(data []byte, addr *net.UDPAddr) (int, error)
}

type Message struct {
	Address *net.UDPAddr
	Data    []byte
	Length  int
}

type messageQueue chan Message

type Server struct {
	config     Config
	conn       *net.UDPConn
	msgQueue   messageQueue
	bufferPool *sync.Pool
	handle     Handler
}

func NewServer(c Config) *Server {
	msgQueue := make(messageQueue, c.MsgQueueSize)

	bufferPool := &sync.Pool{
		New: func() interface{} { return make([]byte, c.MaxBufferSize) },
	}

	return &Server{
		config:     c,
		msgQueue:   msgQueue,
		bufferPool: bufferPool,
	}
}

func (s *Server) OnMessage(h Handler) {
	s.handle = h
}

func (s *Server) StartListening(ctx context.Context, address string) error {
	if s.handle == nil {
		return errors.New("handler function not registered")
	}

	udpAddr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	s.conn = conn

	done := make(chan error, 1)

	for i := 0; i < s.config.NumWorkers; i++ {
		go s.consume()
	}
	for i := 0; i < s.config.NumHandlers; i++ {
		go s.produce(done)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case doneErr := <-done:
		return doneErr
	}
}

func (s *Server) Close() error {
	return s.conn.Close()
}

func (s *Server) WriteToUDP(data []byte, addr *net.UDPAddr) (int, error) {
	return s.conn.WriteToUDP(data, addr)
}

func (s *Server) produce(done chan<- error) {
	for {
		buffer := s.bufferPool.Get().([]byte)
		n, addr, err := s.conn.ReadFromUDP(buffer[0:])
		if err != nil {
			done <- err
		}

		s.msgQueue <- Message{
			Address: addr,
			Data:    buffer[:n],
			Length:  n,
		}
	}
}

func (s *Server) consume() {
	for msg := range s.msgQueue {
		s.process(msg)
	}
}

func (s *Server) process(msg Message) {
	deadline := time.Now().Add(s.config.WriteTimeout)
	if err := s.conn.SetWriteDeadline(deadline); err != nil {
		log.Println(err)
	}

	s.handle(msg, s)
}
