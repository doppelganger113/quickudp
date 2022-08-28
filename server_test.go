//+build !race

package quickudp

import (
	"bytes"
	"context"
	"testing"
)

func TestServer_Running(t *testing.T) {
	addr := "127.0.0.1:1200"
	srv := NewServer(NewConfig())

	srv.OnMessage(func(msg Message, w Writer) {
		_, err := w.WriteToUDP(msg.Data, msg.Address)
		if err != nil {
			t.Fatal(err)
		}
	})

	done := make(chan error, 1)
	go func() {
		done <- srv.StartListening(context.Background(), addr)
	}()

	b, err := NewBenchmark(addr)
	if err != nil {
		t.Fatal(err)
	}

	testData := []struct {
		name string
		data []byte
	}{
		{name: "first", data: []byte("Random blue stuff")},
		{name: "second", data: []byte("Little brown fox jumped over a fence")},
		{name: "third", data: []byte("a")},
	}

	for _, tt := range testData {
		thisTest := tt
		t.Run(thisTest.name, func(t *testing.T) {
			t.Parallel()

			data := []byte("some stuff being sent")
			res, sendErr := b.Send(data)
			if sendErr != nil {
				t.Fatalf("error sending data %v\n", sendErr)
			}

			if bytes.Compare(res, data) != 0 {
				t.Fatal("data sent is not equal to what was received")
			}
		})
	}

	t.Cleanup(func() {
		if err = srv.Close(); err != nil {
			t.Fatalf("error closing the server %v\n", err)
		}
		t.Log(<-done)
	})
}
