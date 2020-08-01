package smux

import (
	"errors"
	"sync"
)

type Allocator interface {
	Get(size int) []byte
	Put(buf []byte) error
}

var defaultAllocator Allocator = NewAllocator()

// Allocator for incoming frames, optimized to prevent overwriting after zeroing
type DefaultAllocator struct {
	buffers []sync.Pool
}

// NewAllocator initiates a []byte allocator for frames less than 65536 bytes,
// the waste(memory fragmentation) of space allocation is guaranteed to be
// no more than 50%.
func NewAllocator() Allocator {
	alloc := new(DefaultAllocator)
	alloc.buffers = make([]sync.Pool, 17) // 1B -> 64K(65536B)
	for k := range alloc.buffers {
		i := k
		alloc.buffers[k].New = func() interface{} {
			return make([]byte, 1<<uint32(i))
		}
	}

	return alloc
}

// Get a []byte from pool with most appropriate cap
func (alloc *DefaultAllocator) Get(size int) []byte {
	if size <= 0 || size > 65536 {
		return nil
	}

	bits := msb(size)
	if size == 1<<bits {
		return alloc.buffers[bits].Get().([]byte)[:size]
	}

	return alloc.buffers[bits+1].Get().([]byte)[:size]
}

// Put returns a []byte to pool for future use,
// which the cap must be exactly 2^n
func (alloc *DefaultAllocator) Put(buf []byte) error {
	bits := msb(cap(buf))
	if cap(buf) == 0 || cap(buf) > 65536 || cap(buf) != 1<<bits {
		return errors.New("allocator Put() incorrect buffer size")
	}
	alloc.buffers[bits].Put(buf)
	return nil
}

// msb return the pos of most significiant bit
// A De Bruijn Sequence is used to find
// the index of the least significant (or first) 1 bit (aka LSB)
// and most significant (or last) 1 bit (aka MSB) set in a given integer.
// https://sites.google.com/site/sydfhd/articles-tutorials/de-bruijn-sequence-generator
// Sean Anderson published bit twiddling hacks containing the Eric Cole's algorithm to find
// the ⌈log2𝑣⌉ of an 𝑁-bit integer 𝑣 in 𝑂(lg(𝑁)) operations with multiply and lookup.
// http://graphics.stanford.edu/~seander/bithacks.html#IntegerLogDeBruijn
// http://supertech.csail.mit.edu/papers/debruijn.pdf
func msb(size int) byte {
	v := uint32(size)

	// 把 Most Significant Bit 以左的位元全部設成 0, 以右全部設為 1。
	// https://hackmd.io/rdTVGkmxSzyTGV9j05qZvw?both
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16

	return debruijinPos[(v*Magic)>>27]
}

var debruijinPos = [32]byte{
	0, 9, 1, 10, 13, 21, 2, 29,
	11, 14, 16, 18, 22, 25, 3, 30,
	8, 12, 20, 28, 15, 17, 24, 7,
	19, 27, 23, 6, 26, 5, 4, 31,
}

// De Bruijn 數列
// 神奇的德布鲁因序列 https://halfrost.com/go_s2_de_bruijn/

const Magic = 0x07C4ACDD // 0x07C4ACDD is 5bit de Bruin Sequence

/*
How to take log2() very fast
https://hackmd.io/@y56/Hk9sTzYWS

ref = https://hackmd.io/rdTVGkmxSzyTGV9j05qZvw?both

My small De Bruijn Test of 130329821:

0x07C4ACDD
00000111110001001010110011011101B
bit 31 - bit 27   00000  0
bit 30 - bit 26   00001  1
bit 29 - bit 25   00011  3
bit 28 - bit 24   00111  7
bit 27 - bit 23   01111 15
bit 26 - bit 22   11111 31
bit 25 - bit 21   11110 30
bit 24 - bit 20   11100 28
bit 23 - bit 19   11000 24
bit 22 - bit 18   10001 17
bit 21 - bit 17   00010  2
bit 20 - bit 16   00100  4
bit 19 - bit 15   01001  9
bit 18 - bit 14   10010 18
bit 17 - bit 13   00101  5
bit 16 - bit 12   01010 10
bit 15 - bit 11   10101 21
bit 14 - bit 10   01011 11
bit 13 - bit  9   10110 22
bit 12 - bit  8   01100 12
bit 11 - bit  7   11001 25
bit 10 - bit  6   10011 19
bit  9 - bit  5   00110  6
bit  8 - bit  4   01101 13
bit  7 - bit  3   11011 27
bit  6 - bit  2   10111 23
bit  5 - bit  1   01110 14
bit  4 - bit  0   11101 29
bit  3 - bit 31   11010 26
bit  2 - bit 30   10100 20
bit  1 - bit 29   01000  8
bit  0 - bit 28   10000 16
*/
