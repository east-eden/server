package encoding2

import (
	"e.coding.net/mmstudio/blade/golib/encoding2/compressor"
)

func GetCompressor(name string) compressor.Compressor {
	return compressor.GetCompressor(name)
}
