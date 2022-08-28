package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"quickudp"
	"time"
)

var (
	address          string
	writeConcurrency int
	readConcurrency  int
	timeout          time.Duration
)

func main() {
	parseArguments()
	fmt.Printf(`Starting benchmark:
	address: %s,
	write concurrency: %d
	read concurrency: %d
`, address, writeConcurrency, readConcurrency)

	bench, err := quickudp.NewBenchmark("localhost:1200")
	if err != nil {
		log.Fatalln("error setting benchmark", err)
	}
	defer func() {
		if closeErr := bench.Close(); closeErr != nil {
			log.Println("error closing benchmark", closeErr)
		}
	}()

	elapsed, err := bench.StartStressTest(context.Background(), writeConcurrency, readConcurrency, timeout)
	if err != nil {
		log.Println("Failed stress test", err)
	}
	bench.PrintConsumedAndProduced()
	bench.CompareResults()
	fmt.Printf("Finished after %s\n", elapsed)
}

func parseArguments() {
	flag.StringVar(&address, "host", "localhost:1200", "Server address to benchmark")
	if address == "" {
		log.Fatalln("host flag should not be empty")
	}
	flag.IntVar(&writeConcurrency, "wc", 1000, "Number of write concurrency")
	if writeConcurrency <= 0 {
		log.Fatalln("c flag should not be 0")
	}
	flag.IntVar(&readConcurrency, "rc", 1000, "Number of read concurrency")
	readTimeout := flag.Int64("t", 15, "Seconds for read timeout")
	if *readTimeout <= 0 || *readTimeout >= 300 {
		log.Fatalln("t flag should have a value between 0 and 300")
	}
	timeout = time.Second * time.Duration(*readTimeout)

	flag.Parse()
}
