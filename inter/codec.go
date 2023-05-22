package inter

type Codec interface {
	Decode(data []byte) (uint32, []byte)
	Encode(data []byte) []byte
	DecodeLength(data []byte) uint32
}
