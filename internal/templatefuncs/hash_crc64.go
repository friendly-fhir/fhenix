package templatefuncs

import "hash/crc64"

type CRC64Module struct{}

func (m *CRC64Module) ISO(s string) uint64 {
	return crc64.Checksum([]byte(s), crc64.MakeTable(crc64.ISO))
}

func (m *CRC64Module) ECMA(s string) uint64 {
	return crc64.Checksum([]byte(s), crc64.MakeTable(crc64.ECMA))
}
