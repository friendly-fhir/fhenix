package templatefuncs

import "hash/fnv"

type FNVModule struct{}

func (m FNVModule) Hash32(s string) uint32 {
	hash := fnv.New32()
	hash.Write([]byte(s))
	return hash.Sum32()
}

func (m FNVModule) Hash32a(s string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	return hash.Sum32()
}

func (m FNVModule) Hash64(s string) uint64 {
	hash := fnv.New64()
	hash.Write([]byte(s))
	return hash.Sum64()
}

func (m FNVModule) Hash64a(s string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(s))
	return hash.Sum64()
}
