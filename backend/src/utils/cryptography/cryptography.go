package cryptography

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
)

func Sha256(str string) string {
	hashValues := []uint32{0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a, 0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19}

	bytes := []byte(str)
	padded := Pad(bytes)
	chunks := Chunks(padded)

	for _, chunk := range chunks {
		ProcessChunk(chunk, &hashValues)
	}

	var output string
	for _, h := range hashValues {
		output += fmt.Sprintf("%08x", h)
	}

	return output
}

func Pad(bytes []byte) []byte {
	originalLen := len(bytes) * 8

	bytes = append(bytes, 0x80)

	for (len(bytes)*8)%512 != 448 {
		bytes = append(bytes, 0x00)
	}

	lenBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(lenBytes, uint64(originalLen))
	bytes = append(bytes, lenBytes...)

	return bytes
}

func Chunks(bytes []byte) [][]byte {
	newLen := int(math.Ceil(float64(len(bytes)) / 64))
	chunks := make([][]byte, newLen)

	for i := 0; i < newLen; i++ {
		start := i * 64
		end := start + 64

		if end > len(bytes) {
			end = len(bytes)
		}

		chunks[i] = bytes[start:end]
	}

	return chunks
}

func ProcessChunk(chunk []byte, H *[]uint32) {
	var K = []uint32{
		0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
		0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
		0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
		0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
		0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
		0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
		0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
		0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
	}

	var W [64]uint32
	for i := 0; i < 16; i++ {
		W[i] = binary.BigEndian.Uint32(chunk[i*4 : (i*4)+4])
	}
	for i := 16; i < 64; i++ {
		W[i] = RotateShiftMix(W[i-2], 17, 19, 10) + W[i-7] + RotateShiftMix(W[i-15], 7, 18, 3) + W[i-16]
	}

	a, b, c, d, e, f, g, h := (*H)[0], (*H)[1], (*H)[2], (*H)[3], (*H)[4], (*H)[5], (*H)[6], (*H)[7]
	for i := 0; i < 64; i++ {
		T1 := h + MajorRotationMix(e, 6, 11, 25) + Ch(e, f, g) + K[i] + W[i]
		T2 := MajorRotationMix(a, 2, 13, 22) + Majority(a, b, c)
		h = g
		g = f
		f = e
		e = d + T1
		d = c
		c = b
		b = a
		a = T1 + T2
	}
	(*H)[0] += a
	(*H)[1] += b
	(*H)[2] += c
	(*H)[3] += d
	(*H)[4] += e
	(*H)[5] += f
	(*H)[6] += g
	(*H)[7] += h
}

func RotateShiftMix(x uint32, rotateA int, rotateB int, shift int) uint32 {
	return (x>>rotateA | x<<(32-rotateA)) ^ (x>>rotateB | x<<(32-rotateB)) ^ (x >> shift)
}

func MajorRotationMix(x uint32, rotateA int, rotateB int, rotateC int) uint32 {
	return (x>>rotateA | x<<(32-rotateA)) ^ (x>>rotateB | x<<(32-rotateB)) ^ (x>>rotateC | x<<(32-rotateC))
}

func Ch(x, y, z uint32) uint32 {
	return (x & y) ^ (^x & z)
}

func Majority(x, y, z uint32) uint32 {
	return (x & y) ^ (x & z) ^ (y & z)
}

func RandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("error while generating a random string: %v", err)
	}

	randomString := base64.URLEncoding.EncodeToString(bytes)

	if len(randomString) > length {
		randomString = randomString[:length]
	}

	return randomString, nil
}

func Base64UrlEncode(input []byte) string {
	return base64.RawURLEncoding.EncodeToString(input)
}

func Base64UrlDecode(input string) ([]byte, error) {
	bytes, err := base64.RawURLEncoding.DecodeString(input)
	if err != nil {
		return nil, fmt.Errorf("error while decoding: %v", err)
	}

	return bytes, nil
}
