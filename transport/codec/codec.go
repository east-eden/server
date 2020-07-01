package codec

// codec define
type CodecType int

const (
	Codec_Protobuf CodecType = iota
	Codec_Json
)

type JsonCodec interface {
	GetName() string
}
