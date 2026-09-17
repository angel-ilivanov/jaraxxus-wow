package protocol

type UpdateObjectResponse struct {
}

func (u UpdateObjectResponse) EncodeBody() []byte {
	// Hardcoded packet from https://gtker.com/wow_messages/docs/smsg_update_object.html?highlight=smsg_update_#example-1-1
	return []byte{
		1, 0, 0, 0, // amount_of_objects: u32
		3,    // [0].Object.update_type: UpdateType CREATE_OBJECT2 (3)
		1, 8, // [0].Object.guid3: PackedGuid
		4,     // [0].Object.object_type: ObjectType PLAYER (4)
		33, 0, // MovementBlock.update_flag: UpdateFlag  SELF| LIVING (33)
		0, 0, 0, 0, 0, 0, // MovementBlock.flags: MovementFlags  NONE (0)
		0, 0, 0, 0, // MovementBlock.timestamp: u32
		205, 215, 11, 198, // Vector4d.x: f32
		53, 126, 4, 195, // Vector4d.y: f32
		249, 15, 167, 66, // Vector4d.z: f32
		0, 0, 0, 0, // Vector4d.orientation: f32
		0, 0, 0, 0, // MovementBlock.fall_time: f32
		0, 0, 128, 63, // MovementBlock.walking_speed: f32
		0, 0, 140, 66, // MovementBlock.running_speed: f32
		0, 0, 144, 64, // MovementBlock.backwards_running_speed: f32
		0, 0, 0, 0, // MovementBlock.swimming_speed: f32
		0, 0, 0, 0, // MovementBlock.backwards_swimming_speed: f32
		0, 0, 0, 0, // MovementBlock.flight_speed: f32
		0, 0, 0, 0, // MovementBlock.backwards_flight_speed: f32
		208, 15, 73, 64, // MovementBlock.turn_rate: f32
		0, 0, 0, 0, // MovementBlock.pitch_rate: f32
		// UpdateMask
		3,          // amount_of_blocks
		7, 0, 0, 0, // Block 0
		0, 0, 128, 0, // Block 1
		24, 0, 0, 0, // Block 2
		8, 0, 0, 0, // Item
		0, 0, 0, 0, // Item
		25, 0, 0, 0, // Item
		1, 0, 0, 0, // Item
		12, 77, 0, 0, // Item
		12, 77, 0, 0, // Item
	}
}

func (u UpdateObjectResponse) Opcode() ServerOpcode {
	return ServerOpcodeUpdateObject
}
