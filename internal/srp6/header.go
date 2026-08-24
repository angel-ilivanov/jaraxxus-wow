package srp6

import (
	hmac2 "crypto/hmac"
	"crypto/rc4"
	"crypto/sha1"
)

var (
	encryptKey = []byte{
		0xCC, 0x98, 0xAE, 0x04, 0xE8, 0x97, 0xEA, 0xCA, 0x12, 0xDD, 0xC0, 0x93, 0x42, 0x91, 0x53, 0x57,
	}
	decryptKey = []byte{
		0xC2, 0xB3, 0x72, 0x3C, 0xC6, 0xAE, 0xD9, 0xB5, 0x34, 0x3C, 0x53, 0xEE, 0x2F, 0x43, 0x67, 0xCE,
	}
)

type HeaderEncryption struct {
	encryptCipher *rc4.Cipher
	decryptCipher *rc4.Cipher
}

func (header *HeaderEncryption) Init(sessionKey []byte) error {
	return header.initKeys(sessionKey, decryptKey, encryptKey)
}

func (header *HeaderEncryption) initKeys(sessionKey, decryptKey, encryptKey []byte) error {
	decryptCipher, err := rc4.NewCipher(createTrafficKey(sessionKey, decryptKey))
	if err != nil {
		return err
	}
	header.decryptCipher = decryptCipher
	discardBytes(header.decryptCipher)
	encryptCipher, err := rc4.NewCipher(createTrafficKey(sessionKey, encryptKey))
	if err != nil {
		return err
	}
	header.encryptCipher = encryptCipher
	discardBytes(header.encryptCipher)
	return nil
}

func createTrafficKey(sessionKey, key []byte) []byte {
	hash := hmac2.New(sha1.New, key)
	hash.Write(sessionKey)
	return hash.Sum(nil)
}

// discardBytes advances the stream by 1024 positions. Done for protection against keystream attack.
// Client does this as well.
func discardBytes(stream *rc4.Cipher) {
	discard := make([]byte, 1024)
	stream.XORKeyStream(discard, discard)
}
