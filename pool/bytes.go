package pool

import (
 "encoding/binary"
 "wl_net/enum"
 "wl_net/inter"
)

type codec struct {
}

func NewCodec() inter.Codec {
 return &codec{}
}

func (c *codec) Decode(data []byte) (uint32, []byte) {
 length := binary.BigEndian.Uint32(data[0:enum.HeadSize])
 dataSize := int(length) - enum.HeadSize
 bytes := make([]byte, dataSize)
 return length, bytes
}

func (c *codec) DecodeLength(data []byte) uint32 {
 length := binary.BigEndian.Uint32(data[0:enum.HeadSize])
 return length
}

func (c *codec) Encode(data []byte) []byte {
 length := uint32(len(data)) + uint32(enum.HeadSize)
 buf := make([]byte, length)
 binary.BigEndian.PutUint32(buf, length)
 buf = append(buf[:enum.HeadSize], data...)
 return buf
}
