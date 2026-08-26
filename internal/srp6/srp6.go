package srp6

import (
	"crypto/rand"
	"crypto/sha1"
	"math/big"
	"strings"
)

const (
	Generator            = 7
	GeneratorLength      = 1
	LargeSafePrimeLength = 32
	saltLength           = 32
	serverKeyLength      = 32
	sessionKeyLength     = 40
	k                    = 3
)

var (
	LargeSafePrimeLittleEndian = []byte{
		0xb7, 0x9b, 0x3e, 0x2a, 0x87, 0x82, 0x3c, 0xab,
		0x8f, 0x5e, 0xbf, 0xbf, 0x8e, 0xb1, 0x01, 0x08,
		0x53, 0x50, 0x06, 0x29, 0x8b, 0x5b, 0xad, 0xbd,
		0x5b, 0x53, 0xe1, 0x89, 0x5e, 0x64, 0x4b, 0x89,
	}
	largeSafePrime  = bytesToBigInt(LargeSafePrimeLittleEndian)
	generatorBigInt = big.NewInt(Generator)
	kBigInt         = big.NewInt(k)
	xorHash         = []byte{
		0xdd, 0x7b, 0xb0, 0x3a, 0x38, 0xac, 0x73, 0x11, 0x3, 0x98,
		0x7c, 0x5a, 0x50, 0x6f, 0xca, 0x96, 0x6c, 0x7b, 0xc2, 0xa7,
	}
)

func GenerateSalt() []byte {
	salt := make([]byte, saltLength)
	_, _ = rand.Read(salt) // never returns an error
	return salt
}

func CalculateWorldServerProof(username string, clientSeed, serverSeed, sessionKey []byte) []byte {
	hash := sha1.New()
	hash.Write([]byte(strings.ToUpper(username)))
	hash.Write(make([]byte, 4))
	hash.Write(clientSeed)
	hash.Write(serverSeed)
	hash.Write(sessionKey)
	return hash.Sum(nil)
}

// CalculatePasswordVerifier returns little endian password verifier
func CalculatePasswordVerifier(username string, password string, salt []byte) []byte {
	x := bytesToBigInt(calculateX(username, password, salt))
	return bigIntToBytes(saltLength, big.NewInt(0).Exp(generatorBigInt, x, largeSafePrime))
}

// x = SHA1( s | SHA1( U | : | p )),
func calculateX(username string, password string, salt []byte) []byte {
	credsHash := sha1.Sum([]byte(strings.ToUpper(username) + ":" + strings.ToUpper(password)))
	concatenated := append(salt[:], credsHash[:]...)
	combinedHash := sha1.Sum(concatenated)
	return combinedHash[:]
}

// CalculateServerPublicKey returns the key as a 32 byte little endian array
func CalculateServerPublicKey(verifier []byte, serverPrivateKey []byte) []byte {
	verifierBigInt := bytesToBigInt(verifier)
	serverPrivateKeyBigInt := bytesToBigInt(serverPrivateKey)

	interim := new(big.Int)
	interim.Mul(kBigInt, verifierBigInt)
	interim.Add(interim, big.NewInt(0).Exp(generatorBigInt, serverPrivateKeyBigInt, largeSafePrime))
	return bigIntToBytes(serverKeyLength, interim.Mod(interim, largeSafePrime))
}

// GenerateServerPrivateKey returns a random 32 byte key
func GenerateServerPrivateKey() []byte {
	key := make([]byte, serverKeyLength)
	_, _ = rand.Read(key)
	return key
}

// serverProof = SHA1(clientPublicKey | clientProof | sessionKey)
func CalculateServerProof(clientPublicKey []byte, clientProof []byte, sessionKey []byte) []byte {
	hash := sha1.New()
	hash.Write(clientPublicKey)
	hash.Write(clientProof)
	hash.Write(sessionKey)
	return hash.Sum(nil)
}

func CalculateExpectedClientProof(username string, sessionKey []byte, clientPublicKey []byte, serverPublicKey []byte, salt []byte) []byte {
	userHash := sha1.Sum([]byte(username))
	hash := sha1.New()
	hash.Write(xorHash)
	hash.Write(userHash[:])
	hash.Write(salt)
	hash.Write(clientPublicKey)
	hash.Write(serverPublicKey)
	hash.Write(sessionKey)
	return hash.Sum(nil)[:]
}

func CalculateServerSessionKey(clientPublicKey []byte, serverPublicKey []byte, passwordVerifier []byte, serverPrivateKey []byte) []byte {
	u := calculateU(clientPublicKey, serverPublicKey)
	sKey := calculateServerSKey(clientPublicKey, passwordVerifier, u, serverPrivateKey)
	return shaInterleave(sKey)
}

// ServerSKey = (clientPublicKey * (verifier^u % largeSafePrime))^serverPrivateKey % largeSafePrime
// Intermediate value for calculating the session key
func calculateServerSKey(clientPublicKey []byte, passwordVerifier []byte, u []byte, serverPrivateKey []byte) []byte {
	clientPublicKeyInt := bytesToBigInt(clientPublicKey)
	passwordVerifierInt := bytesToBigInt(passwordVerifier)
	uInt := bytesToBigInt(u)
	serverPrivateKeyInt := bytesToBigInt(serverPrivateKey)

	result := big.NewInt(0).Exp(passwordVerifierInt, uInt, largeSafePrime)
	result = big.NewInt(0).Mul(clientPublicKeyInt, result)
	result = big.NewInt(0).Exp(result, serverPrivateKeyInt, largeSafePrime)
	return bigIntToBytes(serverKeyLength, result)
}

// u = SHA1( clientPublicKey | serverPublicKey ), intermediate value for calculating the server's S key
func calculateU(clientPublicKey []byte, serverPublicKey []byte) []byte {
	hash := sha1.New()
	hash.Write(clientPublicKey)
	hash.Write(serverPublicKey)
	return hash.Sum(nil)
}

func shaInterleave(sKey []byte) []byte {
	split := splitSKey(sKey) //even length
	var evenBuffer []byte
	var oddBuffer []byte

	for i := 0; i < len(split); i++ {
		if i%2 == 0 {
			evenBuffer = append(evenBuffer, split[i])
		} else {
			oddBuffer = append(oddBuffer, split[i])
		}
	}

	evenHash := sha1.Sum(evenBuffer)
	oddHash := sha1.Sum(oddBuffer)

	sessionKey := make([]byte, sessionKeyLength)
	evenHashIndex := 0
	oddHashIndex := 0
	for i := 0; i < sessionKeyLength; i++ {
		if i%2 == 0 {
			sessionKey[i] = evenHash[evenHashIndex]
			evenHashIndex++
		} else {
			sessionKey[i] = oddHash[oddHashIndex]
			oddHashIndex++
		}
	}
	return sessionKey
}

func splitSKey(sKey []byte) []byte {
	// While the least significant byte is 0, (little endian)
	// remove the 2 least significant bytes.
	// Always keep the length even without
	// trailing zero elements

	for sKey[0] == 0 {
		sKey = sKey[2:]
	}
	return sKey
}

// Adapted from Kangaroux/go-wow-srp6, endian.go:
// https://github.com/Kangaroux/go-wow-srp6/blob/7a61e15fd8d75f4ebe8ac91e07eef140219ef8ff/endian.go
// bytesToBigInt returns a little endian big integer from a big endian byte array.
func bytesToBigInt(data []byte) *big.Int {
	return big.NewInt(0).SetBytes(reverse(data))
}

// bigIntToBytes returns a big endian byte array from a little endian big integer.
func bigIntToBytes(padding int, bi *big.Int) []byte {
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
