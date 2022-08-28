package main

import (
	"context"
	"fmt"
	"log"
	"quickudp"
)

func main() {
	addr := "localhost:1200"
	server := quickudp.NewServer(quickudp.NewConfig())

	server.OnMessage(func(msg quickudp.Message, w quickudp.Writer) {
		if _, err := w.WriteToUDP(msg.Data, msg.Address); err != nil {
			log.Println(err)
		}
	})

	fmt.Printf("UDP server start at %s\n", addr)

	if err := server.StartListening(context.Background(), addr); err != nil {
		log.Println(err)
	}
}
