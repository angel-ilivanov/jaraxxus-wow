package protocol

import (
	"bytes"
	"encoding/binary"
	"testing"
)

var logonChallengeClientPacket = []byte{
	0x00,       // opcode
	0x08,       // protocol version 8
	0x26, 0x00, // size
	'W', 'o', 'W', 0,
	3, 3, 5, // version
	0x34, 0x30, // build 12340 little-endian
	'6', '8', 'x', 0,
	'n', 'i', 'W', 0,
	'S', 'U', 'n', 'e',
	0, 0, 0, 0, // timezone offset
	127, 0, 0, 1, // ip
	8,                                      // account name length
	'J', 'A', 'R', 'A', 'X', 'X', 'U', 'S', // account name
}

var logonProofClientPacket = []byte{
	0x01, // Opcode: CMD_AUTH_LOGON_PROOF
	0xf1, 0x3e, 0xe5, 0xd1, 0x83, 0xc4, 0xc8, 0xa9, 0x50, 0x0e, 0x3f, 0x5a, 0x5d, 0x8a,
	0xee, 0x4e, 0x2e, 0x45, 0xe1, 0xf7, 0xcc, 0x8f, 0x1c, 0xf5, 0xee, 0x8e, 0x11, 0xce,
	0xd3, 0x1d, 0xd7, 0x08, // Client Public Key
	0x6b, 0x1e, 0x48, 0x1b, 0x4d, 0x04, 0xa1, 0x18, 0xd8, 0xf2,
	0xde, 0x5c, 0x59, 0xd5, 0x5c, 0x81, 0x2e, 0x65, 0xec, 0x3e, // Client Proof
	0x4e, 0xf5, 0x2d, 0xe1,
	0x80, 0x5e, 0x1a, 0x67, 0x15, 0xec, 0xc8, 0x41, 0xee, 0xb8, 0x90, 0x8a, 0x58, 0xbb,
	0x00, 0xd0, // CRC Hash
	0x00, // Num keys: 0
	0x00, // Two factor enabled: false
}

func TestDecodeRequestEmptyPacket(t *testing.T) {
	_, err := DecodeRequest(bytes.NewReader(make([]byte, 0)))
	if err == nil {
		t.Fatalf("Decoder accepts empty packet, expected to fail")
	}
}

func TestDecodeRequestUnknownOpcode(t *testing.T) {
	_, err := DecodeRequest(bytes.NewReader([]byte{0xFF}))
	if err == nil {
		t.Fatalf("Decoder accepts unknown opcode, expected to fail")
	}
}

func TestDecodeUsername(t *testing.T) {
	result, err := DecodeRequest(bytes.NewReader(logonChallengeClientPacket))
	if err != nil {
		t.Fatalf("ReadClientMessage() error = %v", err)
	}
	request, ok := result.(LogonChallengeRequest)
	if !ok {
		t.Fatalf("Expected a LogonChallengeClientRequest, but got %T", result)
	}
	if request.AccountName != "JARAXXUS" {
		t.Fatalf("Expected JARAXXUS but got %s", request.AccountName)
	}
}

func TestDecodeIp(t *testing.T) {
	result, err := DecodeRequest(bytes.NewReader(logonChallengeClientPacket))
	if err != nil {
		t.Fatalf("ReadClientMessage() error = %v", err)
	}
	request, ok := result.(LogonChallengeRequest)
	if !ok {
		t.Fatalf("Expected a LogonChallengeClientRequest, but got %T", result)
	}
	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, request.IP)
	if request.IP != 2130706433 {
		t.Fatalf("Expected Ip [127 0 0 1] but got %v", ipBytes)
	}
}

func TestDecodeClientProof(t *testing.T) {
	result, err := DecodeRequest(bytes.NewReader(logonProofClientPacket))
	if err != nil {
		t.Fatalf("ReadClientMessage() error = %v", err)
	}
	request, ok := result.(LogonProofRequest)
	if !ok {
		t.Fatalf("Expected a LogonProofClientRequest, but got %T", result)
	}
	expected := []byte{
		0x6b, 0x1e, 0x48, 0x1b, 0x4d, 0x04, 0xa1, 0x18, 0xd8, 0xf2,
		0xde, 0x5c, 0x59, 0xd5, 0x5c, 0x81, 0x2e, 0x65, 0xec, 0x3e,
	}
	if bytes.Compare(request.ClientProof, expected) != 0 {
		t.Fatalf("Server incorrectly decodes client proof, expected %v, but got %v", expected, request.ClientProof)
	}
}

func TestDecodeClientPublicKey(t *testing.T) {
	result, err := DecodeRequest(bytes.NewReader(logonProofClientPacket))
	if err != nil {
		t.Fatalf("ReadClientMessage() error = %v", err)
	}
	request, ok := result.(LogonProofRequest)
	if !ok {
		t.Fatalf("Expected a LogonProofClientRequest, but got %T", result)
	}
	expected := []byte{
		0xf1, 0x3e, 0xe5, 0xd1, 0x83, 0xc4, 0xc8, 0xa9, 0x50, 0x0e, 0x3f, 0x5a, 0x5d, 0x8a,
		0xee, 0x4e, 0x2e, 0x45, 0xe1, 0xf7, 0xcc, 0x8f, 0x1c, 0xf5, 0xee, 0x8e, 0x11, 0xce,
		0xd3, 0x1d, 0xd7, 0x08,
	}
	if bytes.Compare(request.ClientPublicKey, expected) != 0 {
		t.Fatalf("Server incorrectly decodes client proof, expected %v, but got %v", expected, request.ClientPublicKey)
	}
}
