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

func TestReadClientMessageEmptyPacket(t *testing.T) {
	_, err := ReadClientMessage(bytes.NewReader(make([]byte, 0)))
	if err == nil {
		t.Fatalf("Decoder accepts empty packet, expected to fail")
	}
}

func TestReadClientMessageUnknownOpcode(t *testing.T) {
	_, err := ReadClientMessage(bytes.NewReader([]byte{0xFF}))
	if err == nil {
		t.Fatalf("Decoder accepts unknown opcode, expected to fail")
	}
}

func TestDecodeUsername(t *testing.T) {
	result, err := ReadClientMessage(bytes.NewReader(logonChallengeClientPacket))
	if err != nil {
		t.Fatalf("ReadClientMessage() error = %v", err)
	}
	request, ok := result.(LogonChallengeClientMessage)
	if !ok {
		t.Fatalf("Expected a LogonChallengeClientRequest, but got %T", result)
	}
	if request.AccountName != "JARAXXUS" {
		t.Fatalf("Expected JARAXXUS but got %s", request.AccountName)
	}
}

func TestDecodeIp(t *testing.T) {
	result, err := ReadClientMessage(bytes.NewReader(logonChallengeClientPacket))
	if err != nil {
		t.Fatalf("ReadClientMessage() error = %v", err)
	}
	request, ok := result.(LogonChallengeClientMessage)
	if !ok {
		t.Fatalf("Expected a LogonChallengeClientRequest, but got %T", result)
	}
	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, request.Ip)
	if request.Ip != 2130706433 {
		t.Fatalf("Expected Ip [127 0 0 1] but got %v", ipBytes)
	}
}
