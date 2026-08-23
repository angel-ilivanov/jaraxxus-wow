package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type logonProofWire struct {
	Opcode          uint8
	ClientPublicKey [32]uint8
	ClientProof     [20]uint8
	CrcHash         [20]uint8
	NumKeys         uint8
	SecurityFlags   uint8
}

type LogonProofRequest struct {
	ClientPublicKey []byte
	ClientProof     []byte
}

func (l LogonProofRequest) isRequest() {}

func DecodeLogonProofRequest(reader io.Reader) (LogonProofRequest, error) {
	packetBytes, err := readLogonProofPacket(reader)
	if err != nil {
		return LogonProofRequest{}, fmt.Errorf("error decoding logon proof request: %w", err)
	}
	wire, err := decodeLogonProofWire(packetBytes)
	if err != nil {
		return LogonProofRequest{}, fmt.Errorf("error decoding logon proof request: %w", err)
	}
	return newLogonProofRequest(wire), nil
}

func readLogonProofPacket(reader io.Reader) ([]byte, error) {
	packet := make([]byte, 74)
	_, err := io.ReadFull(reader, packet)
	if err != nil {
		return nil, fmt.Errorf("error reading logon proof packet: %w", err)
	}
	packet = append([]byte{byte(CmdAuthLogonProof)}, packet...) //prepend opcode
	return packet, nil
}

func decodeLogonProofWire(packetBytes []byte) (logonProofWire, error) {
	var parsed logonProofWire
	reader := bytes.NewReader(packetBytes)
	err := binary.Read(reader, binary.LittleEndian, &parsed)
	if err != nil {
		return logonProofWire{}, fmt.Errorf("error decoding packet bytes into wire object: %w", err)
	}
	return parsed, nil
}

func newLogonProofRequest(wire logonProofWire) LogonProofRequest {
	return LogonProofRequest{
		ClientPublicKey: wire.ClientPublicKey[:],
		ClientProof:     wire.ClientProof[:],
	}
}
