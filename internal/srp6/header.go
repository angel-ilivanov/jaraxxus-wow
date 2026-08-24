package srp6

import (
	hmac2 "crypto/hmac"
	"crypto/rc4"
	"crypto/sha1"
	"errors"
)

var (
	serverToClientSeed = []byte{
		0xCC, 0x98, 0xAE, 0x04, 0xE8, 0x97, 0xEA, 0xCA, 0x12, 0xDD, 0xC0, 0x93, 0x42, 0x91, 0x53, 0x57,
	}
	clientToServerSeed = []byte{
		0xC2, 0xB3, 0x72, 0x3C, 0xC6, 0xAE, 0xD9, 0xB5, 0x34, 0x3C, 0x53, 0xEE, 0x2F, 0x43, 0x67, 0xCE,
	}
)

type HeaderEncryption struct {
	sendCipher    *rc4.Cipher
	receiveCipher *rc4.Cipher
}

var ErrUninitializedCipher = errors.New("rc4 cipher cannot be used before initialization")

// EncryptHeader encrypts header in place
func (h *HeaderEncryption) EncryptHeader(header []byte) error {
	if h.sendCipher == nil {
		return ErrUninitializedCipher
	}
	h.sendCipher.XORKeyStream(header, header)
	return nil
}

func (h *HeaderEncryption) DecryptHeader(header []byte) error {
	if h.receiveCipher == nil {
		return ErrUninitializedCipher
	}
	h.receiveCipher.XORKeyStream(header, header)
	return nil
}

func (h *HeaderEncryption) Init(sessionKey []byte) error {
	return h.initKeys(sessionKey, clientToServerSeed, serverToClientSeed)
}

func (h *HeaderEncryption) initKeys(sessionKey, receiveSeed, sendSeed []byte) error {
	sendStream, err := newStream(sessionKey, sendSeed)
	if err != nil {
		return err
	}
	receiveStream, err := newStream(sessionKey, receiveSeed)
	if err != nil {
		return err
	}
	h.sendCipher = sendStream
	h.receiveCipher = receiveStream
	return nil
}

func newStream(sessionKey, key []byte) (*rc4.Cipher, error) {
	stream, err := rc4.NewCipher(createTrafficKey(sessionKey, key))
	if err != nil {
		return nil, err
	}
	discardBytes(stream)
	return stream, nil
}

func createTrafficKey(sessionKey, key []byte) []byte {
	hash := hmac2.New(sha1.New, key)
	hash.Write(sessionKey)
	return hash.Sum(nil)
}

// discardBytes discards the first 1024 RC4 keystream bytes as required by Wrath.
func discardBytes(stream *rc4.Cipher) {
	discard := make([]byte, 1024)
	stream.XORKeyStream(discard, discard)
}
