package quickudp

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Benchmark struct {
	sync.RWMutex
	produced map[int]uint
	consumed map[int]uint
	config   BenchConfig
	conn     *net.UDPConn
	buffPool *sync.Pool
}

func NewBenchmark(address string, options ...BenchOption) (*Benchmark, error) {
	benchmark := &Benchmark{
		produced: map[int]uint{},
		consumed: map[int]uint{},
	}
	benchmark.config = NewBenchConfig(options...)

	udpAddr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}
	benchmark.conn = conn

	benchmark.buffPool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, benchmark.config.MaxBufferSize)
		},
	}

	return benchmark, nil
}

func (b *Benchmark) Write(data []byte) (int, error) {
	return b.conn.Write(data)
}

func (b *Benchmark) StartStressTest(
	ctx context.Context,
	writeConcurrency int,
	readConcurrency int,
	readDeadline time.Duration,
) (elapsedTime time.Duration, err error) {
	if err = b.conn.SetReadDeadline(time.Now().Add(readDeadline)); err != nil {
		return 0, err
	}

	done := make(chan error, 1)
	start := time.Now()

	go func() {
		var wg sync.WaitGroup
		wg.Add(readConcurrency + writeConcurrency)

		for i := 0; i < readConcurrency; i++ {
			go func() {
				defer wg.Done()
				if consumeErr := b.consume(); consumeErr != nil {
					done <- fmt.Errorf("failed consuming: %v", consumeErr)
				}
			}()
		}

		for i := 0; i < writeConcurrency; i++ {
			go func() {
				defer wg.Done()
				if produceErr := b.produce(); produceErr != nil {
					done <- fmt.Errorf("failed consuming: %v", produceErr)
				}
			}()
		}

		wg.Wait()
		done <- nil
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case err = <-done:
		if err != nil {
			return 0, err
		}
		elapsedTime = time.Since(start)

		return elapsedTime, nil
	}
}

func (b *Benchmark) Close() error {
	return b.conn.Close()
}

func (b *Benchmark) produce() error {
	value := b.config.DummyData[rand.Intn(len(b.config.DummyData))]

	_, sendErr := b.Write(value)
	if sendErr != nil {
		return sendErr
	}

	b.incrementProduced(value)

	return nil
}

func (b *Benchmark) incrementProduced(data []byte) {
	b.Lock()
	b.produced[len(data)] = b.produced[len(data)] + 1
	b.Unlock()
}

func (b *Benchmark) incrementConsumed(data []byte) {
	b.Lock()
	b.consumed[len(data)] = b.consumed[len(data)] + 1
	b.Unlock()
}

func (b *Benchmark) consume() error {
	buff := b.buffPool.Get().([]byte)
	j, err := b.conn.Read(buff[0:])
	if err != nil {
		return err
	}

	data := buff[0:j]
	b.incrementConsumed(data)

	return nil
}

func (b *Benchmark) PrintConsumedAndProduced() {
	fmt.Println("---------------------------------------")
	fmt.Println("Produced")
	var count uint
	for key, value := range b.produced {
		fmt.Printf("\tKey: %d Count: %d\n", key, value)
		count += value
	}
	fmt.Printf("Total produced: %d\n", count)

	fmt.Println("---------------------------------------")
	fmt.Println("Consumed")
	count = 0
	for key, value := range b.consumed {
		fmt.Printf("\tKey: %d Count: %d\n", key, value)
		count += value
	}
	fmt.Printf("Total consumed: %d\n", count)
	fmt.Println("---------------------------------------")
}

func (b *Benchmark) CompareResults() {
	for length, num := range b.produced {
		dummyValue := getDummyDataByLength(b.config.DummyData, length)
		count, ok := b.consumed[length]
		if !ok {
			log.Fatalf("Missing in consumed: %s\n", dummyValue)
		}
		if count != num {
			log.Fatalf("Data %s has different count values: expected %d got %d\n", dummyValue, count, num)
		}
	}
}
