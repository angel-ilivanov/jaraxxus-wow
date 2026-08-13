package srp6

import (
	"crypto/rand"
	"crypto/sha1"
	"math/big"
	"strings"
)

const generator = 7

var largeSafePrimeLittleEndian = []byte{
	0xb7, 0x9b, 0x3e, 0x2a, 0x87, 0x82, 0x3c, 0xab,
	0x8f, 0x5e, 0xbf, 0xbf, 0x8e, 0xb1, 0x01, 0x08,
	0x53, 0x50, 0x06, 0x29, 0x8b, 0x5b, 0xad, 0xbd,
	0x5b, 0x53, 0xe1, 0x89, 0x5e, 0x64, 0x4b, 0x89,
}
var n = bytesToInt(largeSafePrimeLittleEndian)
var g = big.NewInt(generator)

func generateSalt() []byte {
	salt := make([]byte, 32)
	rand.Read(salt) // never returns an error
	return salt
}

// returns little endian password verifier
func calculatePasswordVerifier(username string, password string, salt []byte) []byte {
	x := bytesToInt(calculateX(username, password, salt))
	return intToBytes(32, big.NewInt(0).Exp(g, x, n))
}

// salt and X are in little endian
func calculateX(username string, password string, salt []byte) []byte {
	interim := sha1.Sum([]byte(strings.ToUpper(username) + ":" + strings.ToUpper(password)))
	concatenated := append(salt[:], interim[:]...)
	hash := sha1.Sum(concatenated)
	return hash[:]
}

// Adapted from Kangaroux/go-wow-srp6, endian.go:
// https://github.com/Kangaroux/go-wow-srp6/blob/7a61e15fd8d75f4ebe8ac91e07eef140219ef8ff/endian.go
// bytesToInt returns a little endian big integer from a big endian byte array.
func bytesToInt(data []byte) *big.Int {
	return big.NewInt(0).SetBytes(reverse(data))
}

// intToBytes returns a big endian byte array from a little endian big integer.
func intToBytes(padding int, bi *big.Int) []byte {
	return reverse(pad(padding, bi.Bytes()))
}

func pad(length int, data []byte) []byte {
	dataLen := len(data)
	if dataLen >= length {
		return data
	}
	ret := make([]byte, length)
	copy(ret[length-dataLen:], data)
	return ret
}

// reverse returns a copy of data in reverse order.
func reverse(data []byte) []byte {
	n := len(data)
	newData := make([]byte, n)
	for i := 0; i < n; i++ {
		newData[i] = data[n-i-1]
	}
	return newData
}
