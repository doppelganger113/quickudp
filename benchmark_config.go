package quickudp

var standardDummyData = [][]byte{
	[]byte("test"),
	[]byte("a little bit longer byte slice with a chance to detect an error of some sorts, we'll see"),
	[]byte("little brown fox jumped over a fence"),
	[]byte("a"),
}

func getDummyDataByLength(data [][]byte, length int) []byte {
	for _, value := range data {
		if len(value) == length {
			return value
		}
	}

	return nil
}

type BenchConfig struct {
	DummyData     [][]byte
	MaxBufferSize int
}

var defaultBenchConfig = BenchConfig{
	DummyData:     standardDummyData,
	MaxBufferSize: 512,
}

func NewBenchConfig(options ...BenchOption) BenchConfig {
	c := defaultBenchConfig

	for _, o := range options {
		o(&c)
	}

	return c
}

type BenchOption func(b *BenchConfig)

func WithBenchDummyData(data [][]byte) BenchOption {
	return func(b *BenchConfig) {
		b.DummyData = data
	}
}

func WithBenchMaxBufferSize(size int) BenchOption {
	return func(b *BenchConfig) {
		b.MaxBufferSize = size
	}
}
