package bitcask

import (
	"hash/crc32"

	pb "github.com/prologic/bitcask/proto"
)

func NewEntry(key string, value []byte) pb.Entry {
	crc := crc32.ChecksumIEEE(value)

	return pb.Entry{
		CRC:   crc,
		Key:   key,
		Value: value,
	}
}
