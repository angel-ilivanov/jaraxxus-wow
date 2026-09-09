package characterstore

type Character struct {
	GUID       uint32
	Name       string
	Race       Race
	Class      Class
	Gender     Gender
	Appearance Appearance
	State      State
}

type Appearance struct {
	Skin        uint8
	Face        uint8
	HairStyle   uint8
	HairColor   uint8
	FacialStyle uint8
}

type State struct {
	Level       uint8
	MapID       uint32
	ZoneID      uint32
	PositionX   float32
	PositionY   float32
	PositionZ   float32
	Orientation float32
}

type Gender uint8

const (
	GenderMale   Gender = 0
	GenderFemale Gender = 1
)

type Race uint8

const (
	RaceHuman    Race = 1
	RaceOrc      Race = 2
	RaceDwarf    Race = 3
	RaceNightElf Race = 4
	RaceUndead   Race = 5
	RaceTauren   Race = 6
	RaceGnome    Race = 7
	RaceTroll    Race = 8
	RaceBloodElf Race = 10
	RaceDraenei  Race = 11
)

type Class uint8

const (
	ClassWarrior     Class = 1
	ClassPaladin     Class = 2
	ClassHunter      Class = 3
	ClassRogue       Class = 4
	ClassPriest      Class = 5
	ClassDeathKnight Class = 6
	ClassShaman      Class = 7
	ClassMage        Class = 8
	ClassWarlock     Class = 9
	ClassDruid       Class = 11
)
