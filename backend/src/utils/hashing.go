package main

import (
	"encoding/binary"
	"fmt"
)

func SHA256(str string) {
	bytes := []byte(str)
	bytesPadded := pad(bytes)

	fmt.Printf("Bytes: %08b\n", bytes)
	fmt.Printf("Padded Bytes: %08b\n", bytesPadded)
}

func pad(bytes []byte) []byte {
	originalLen := len(bytes) * 8

	// A single [1] bit must be appended.
	bytes = append(bytes, 0x80)

	// [0] bits must be added until their length is equal 448 when % 512.
	// This means that there are 64 bits left for storing the big-endian integer.
	for (len(bytes)*8)%512 != 448 {
		fmt.Printf("%08b\n", bytes)
		bytes = append(bytes, 0x00)
	}

	// The free 64 bit space is for the big-endian integer representing the
	// length of the string.
	lenBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(lenBytes, uint64(originalLen))
	bytes = append(bytes, lenBytes...)

	// We're left with a byte array to which was added a single [1] bit.
	// After that it was filled with [0] bits until there were only 64 bits left.
	// The 64 bits were assigned to the big-endian representation of the string's length.
	return bytes
}

func main() {
	SHA256("hello world")
}
