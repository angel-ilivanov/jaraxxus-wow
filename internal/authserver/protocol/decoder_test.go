package protocol

import (
	"bytes"
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

func TestReadClientMessageReturnsLogonChallengeRequest(t *testing.T) {
	result, err := ReadClientMessage(bytes.NewReader(logonChallengeClientPacket))
	if err != nil {
		t.Fatalf("ReadClientMessage() error = %v", err)
	}
	_, ok := result.(LogonChallengeClientRequest)
	if !ok {
		t.Fatalf("Expected a LogonChallengeClientRequest, but got %T", result)
	}
}

func TestDecodeUsername(t *testing.T) {
	result, err := ReadClientMessage(bytes.NewReader(logonChallengeClientPacket))
	if err != nil {
		t.Fatalf("ReadClientMessage() error = %v", err)
	}
	name := result.(LogonChallengeClientRequest).AccountName
	if name != "JARAXXUS" {
		t.Fatalf("Expected JARAXXUS but got %s", name)
	}
}

func TestDecodeIp(t *testing.T) {
	result, err := ReadClientMessage(bytes.NewReader(logonChallengeClientPacket))
	if err != nil {
		t.Fatalf("ReadClientMessage() error = %v", err)
	}
	ip := result.(LogonChallengeClientRequest).Ip
	if ip != 2130706433 {
		t.Fatalf("Expected Ip 2130706433 but got %d", ip)
	}
}
