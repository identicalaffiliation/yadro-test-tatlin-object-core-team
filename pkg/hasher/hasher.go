package hasher

import "hash/fnv"

type hasher struct{}

func NewHasher() *hasher { return &hasher{} }

func (hs *hasher) HashLine32(line string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(line)) //nolint: errcheck
	return h.Sum32()
}
