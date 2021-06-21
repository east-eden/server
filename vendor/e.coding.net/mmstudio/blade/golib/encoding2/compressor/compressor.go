package compressor

import "io"

type WriteFlushCloser interface {
	io.WriteCloser
	Flush() error
}

type Compressor interface {
	Compress(w io.Writer) (WriteFlushCloser, error)
	Decompress(r io.Reader) (io.ReadCloser, error)
	Name() string
}

var registeredCompressor = make(map[string]Compressor)

func RegisterCompressor(c Compressor) {
	registeredCompressor[c.Name()] = c
}
func GetCompressor(name string) Compressor {
	return registeredCompressor[name]
}
