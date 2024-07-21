package templatefuncs

import "hash/crc32"

type CRC32Module struct{}

func (m *CRC32Module) IEEE(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}

func (m *CRC32Module) Castagnoli(s string) uint32 {
	table := crc32.MakeTable(crc32.Castagnoli)
	return crc32.Checksum([]byte(s), table)
}

func (m *CRC32Module) Koopman(s string) uint32 {
	table := crc32.MakeTable(crc32.Koopman)
	return crc32.Checksum([]byte(s), table)
}
