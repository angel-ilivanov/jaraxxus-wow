package protocol

import (
	"bytes"
	"encoding/binary"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/characterstore"
)

type CharEnumRequest struct{}

func (c CharEnumRequest) isClientMessage() {}

func (c CharEnumRequest) Opcode() ClientOpcode {
	return ClientOpcodeCharEnum
}

type CharEnumResponse struct {
	Characters []characterstore.Character
}

func (c CharEnumResponse) EncodeBody() []byte {
	var buf bytes.Buffer

	// Characters Header:
	buf.WriteByte(byte(len(c.Characters)))
	if len(c.Characters) == 0 {
		return buf.Bytes()
	}
	// characters...
	for _, char := range c.Characters {
		buf.Write(encodeCharacter(char))
	}
	return buf.Bytes()
}

func encodeCharacter(character characterstore.Character) []byte {
	var buf bytes.Buffer

	guidBuffer := make([]byte, 0, 8)
	guidBuffer = binary.LittleEndian.AppendUint64(guidBuffer, character.GUID)
	buf.Write(guidBuffer)

	buf.WriteString(character.Name)
	buf.WriteByte(0) // CString terminator

	buf.WriteByte(byte(character.Race))
	buf.WriteByte(byte(character.Class))
	buf.WriteByte(byte(character.Gender))
	buf.WriteByte(character.Appearance.Skin)
	buf.WriteByte(character.Appearance.Face)
	buf.WriteByte(character.Appearance.HairStyle)
	buf.WriteByte(character.Appearance.HairColor)
	buf.WriteByte(character.Appearance.FacialStyle)

	buf.WriteByte(character.State.Level)

	areaBuf := make([]byte, 0, 4)
	areaBuf = binary.LittleEndian.AppendUint32(areaBuf, character.State.ZoneID)
	buf.Write(areaBuf)

	mapBuf := make([]byte, 0, 4)
	mapBuf = binary.LittleEndian.AppendUint32(mapBuf, character.State.MapID)
	buf.Write(mapBuf)

	writeFloat32LE(&buf, character.State.PositionX)
	writeFloat32LE(&buf, character.State.PositionY)
	writeFloat32LE(&buf, character.State.PositionZ)

	padding32Bytes := make([]byte, 4)
	buf.Write(padding32Bytes) // guild id, not implemented
	buf.Write(padding32Bytes) // TODO character flags (none, hideHelm, hideCloak, ghost)
	buf.Write(padding32Bytes) // recustomization_flags, not implemented
	buf.WriteByte(0)          // character's first log in (false), hardcoded disabled tutorials
	buf.Write(padding32Bytes) // character_pet_display_id, not implemented
	buf.Write(padding32Bytes) // character_pet_level, not implemented
	buf.Write(padding32Bytes) // character_pet_family, not implemented

	buf.Write(make([]byte, 23*9)) //gear slots, hardcoded to empty

	return buf.Bytes()
}

func (c CharEnumResponse) Opcode() ServerOpcode {
	return ServerOpcodeCharEnum
}
