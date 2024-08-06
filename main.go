package main

import (
	"fmt"
	"hash"

	"github.com/google/uuid"
	"github.com/spaolacci/murmur3"
)

var mhasher hash.Hash32

func init() {
	//mhasher = murmur3.New32WithSeed(uint32(time.Now().UnixNano()))
	mhasher = murmur3.New32WithSeed(uint32(10))
}
func murmurhash(key string, size int32) int32 {
	mhasher.Write([]byte(key))
	result := mhasher.Sum32() % uint32(size)
	mhasher.Reset()
	return int32(result)
}

type BloomFilter struct {
	filter []byte
	size   int32
}

func NewBloomFilter(size int32) *BloomFilter {
	return &BloomFilter{
		filter: make([]byte, size),
		size:   size,
	}

}

func (b *BloomFilter) Add(key string) {
	idx := murmurhash(key, b.size)
	b.filter[idx/8] = b.filter[idx/8] | (1 << (idx % 8))
	//fmt.Println("Wrote", key, "to", idx)
}

func (b *BloomFilter) Exists(key string) bool {
	idx := murmurhash(key, b.size)
	return b.filter[idx/8]&(1<<(idx%8)) != 0
}

func (b *BloomFilter) Print() {
	for _, v := range b.filter {
		fmt.Println(v)
	}
}

func main() {

	dataset := make([]string, 0)
	dataset_exists := make(map[string]bool)
	dataset_doesnotexists := make(map[string]bool)

	for i := 0; i < 500; i++ {
		u := uuid.New()
		dataset = append(dataset, u.String())
		dataset_exists[u.String()] = true
	}

	for i := 0; i < 500; i++ {
		u := uuid.New()
		dataset = append(dataset, u.String())
		dataset_doesnotexists[u.String()] = true
	}

	bloom := NewBloomFilter(1600)

	for key := range dataset_exists {
		bloom.Add(key)
	}
	falsePositive := 0
	for _, key := range dataset {
		exists := bloom.Exists(key)
		if exists {
			if _, ok := dataset_doesnotexists[key]; ok {
				falsePositive++
			}
		}
	}
	fmt.Println("False Positive", 100*float64(falsePositive)/float64(len(dataset_doesnotexists)))
}
