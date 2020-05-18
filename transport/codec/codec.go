package codec

// codec define
const (
	Codec_Protobuf = iota
	Codec_Json
)

type JsonCodec interface {
	GetName() string
}
