package protocol

import (
	"bytes"
	"encoding/binary"
	"math"
)

func writeFloat32LE(buf *bytes.Buffer, value float32) {
	var data [4]byte
	binary.LittleEndian.PutUint32(data[:], math.Float32bits(value))
	buf.Write(data[:])
}

func writeUint32LE(buf *bytes.Buffer, value uint32) {
	data := make([]byte, 0, 4)
	binary.LittleEndian.AppendUint32(data, value)
	buf.Write(data)
}

func writeUint64LE(buf *bytes.Buffer, value uint64) {
	data := make([]byte, 0, 8)
	binary.LittleEndian.AppendUint64(data, value)
	buf.Write(data)
}
