package quickudp

import "time"

type Config struct {
	MaxBufferSize int
	MsgQueueSize  int
	NumWorkers    int
	NumHandlers   int
	WriteTimeout  time.Duration
	MemTickPeriod time.Duration
}

var defaultConfig = Config{
	MaxBufferSize: 512,
	MsgQueueSize:  100_000,
	NumWorkers:    5_000,
	NumHandlers:   100,
	WriteTimeout:  1 * time.Second,
	MemTickPeriod: 10 * time.Second,
}

type Option func(c *Config)

func NewConfig(options ...Option) Config {
	c := defaultConfig
	for _, o := range options {
		o(&c)
	}
	return c
}

func WithMaxBufferSize(size int) Option {
	return func(c *Config) {
		c.MaxBufferSize = size
	}
}

func WithMsgQueueSize(size int) Option {
	return func(c *Config) {
		c.MsgQueueSize = size
	}
}

func WithNumWorkers(numWorkers int) Option {
	return func(c *Config) {
		c.NumWorkers = numWorkers
	}
}

func WithNumHandlers(numHandlers int) Option {
	return func(c *Config) {
		c.NumHandlers = numHandlers
	}
}
