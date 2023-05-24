package snowflake

import (
	"github.com/bwmarrin/snowflake"
)

var (
	_snowflakeGlobal *snowflake.Node
	_snowflakeRegion *snowflake.Node
)

func SetSnowflakeGlobalNodeId(id int64) (err error) {
	_snowflakeGlobal, err = snowflake.NewNode(id)
	return
}

func GenSnowflakeGlobalNodeId() int64 {
	return _snowflakeGlobal.Generate().Int64()
}

func SetSnowflakeRegionNodeId(nodeId int64) (err error) {
	_snowflakeRegion, err = snowflake.NewNode(nodeId)
	return
}

func GenSnowflakeRegionNodeId() int64 {
	return _snowflakeRegion.Generate().Int64()
}

func Int64ToBytes(id int64) []byte {
	return []byte{
		byte(id >> 56),
		byte(id >> 48),
		byte(id >> 40),
		byte(id >> 32),
		byte(id >> 24),
		byte(id >> 16),
		byte(id >> 8),
		byte(id),
	}
}

func BytesToInt64(bytes []byte) int64 {
	return int64(bytes[0])<<56 |
		int64(bytes[1])<<48 |
		int64(bytes[2])<<40 |
		int64(bytes[3])<<32 |
		int64(bytes[4])<<24 |
		int64(bytes[5])<<16 |
		int64(bytes[6])<<8 |
		int64(bytes[7])
}

func Int64ToStr(id int64) string {
	return string(Int64ToBytes(id))
}

func Int64sToStrSlice(ids ...int64) []string {
	strings := make([]string, len(ids))
	for i, id := range ids {
		strings[i] = Int64ToStr(id)
	}
	return strings
}

func StrToInt64(id string) int64 {
	return BytesToInt64([]byte(id))
}

func Uint64ToStr(id uint64) string {
	return string([]byte{
		byte(id >> 56),
		byte(id >> 48),
		byte(id >> 40),
		byte(id >> 32),
		byte(id >> 24),
		byte(id >> 16),
		byte(id >> 8),
		byte(id),
	})
}

func StrToUint64(id string) uint64 {
	bytes := []byte(id)
	return uint64(bytes[0])<<56 |
		uint64(bytes[1])<<48 |
		uint64(bytes[2])<<40 |
		uint64(bytes[3])<<32 |
		uint64(bytes[4])<<24 |
		uint64(bytes[5])<<16 |
		uint64(bytes[6])<<8 |
		uint64(bytes[7])
}

func Uint32ToStr(id uint32) string {
	return string([]byte{
		byte(id >> 24),
		byte(id >> 16),
		byte(id >> 8),
		byte(id),
	})
}

func StrToUint32(id string) uint32 {
	bytes := []byte(id)
	return uint32(bytes[0])<<24 |
		uint32(bytes[1])<<16 |
		uint32(bytes[2])<<8 |
		uint32(bytes[3])
}

func Uint16ToStr(id uint16) string {
	return string([]byte{
		byte(id >> 8),
		byte(id),
	})
}

func StrToUint16(id string) uint16 {
	bytes := []byte(id)
	return uint16(bytes[0])<<8 |
		uint16(bytes[1])
}

func MergeUInt32(prevId, lastId uint16) uint32 {
	return uint32(prevId)<<16 | uint32(lastId)
}

func SplitUInt32(code uint32) (prevId, lastId uint16) {
	prevId = uint16(code >> 16)
	lastId = uint16(code)
	return
}
