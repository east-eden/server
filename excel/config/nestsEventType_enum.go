package config

// 音巢挑战-拳皇/NestsEventType
type NestsEventType int

const (
	NestsEventType_BEGIN       NestsEventType = 0  //枚举开始
	NestsEventType_None        NestsEventType = 0  //N
	NestsEventType_Award       NestsEventType = 1  //A
	NestsEventType_Little      NestsEventType = 2  //B
	NestsEventType_Elite       NestsEventType = 3  //C
	NestsEventType_Boss        NestsEventType = 4  //D
	NestsEventType_SecretRoom  NestsEventType = 5  //E
	NestsEventType_Obstacle    NestsEventType = 6  //Z
	NestsEventType_Trap        NestsEventType = 7  //T
	NestsEventType_Random      NestsEventType = 8  //R
	NestsEventType_RedEnvelope NestsEventType = 9  //RE
	NestsEventType_TrapDebuff  NestsEventType = 10 //TD
	NestsEventType_END         NestsEventType = 11 //枚举结束
)

func (e NestsEventType) String() string {

	if e == NestsEventType_None {
		return "None"
	}
	if e == NestsEventType_Award {
		return "Award"
	}
	if e == NestsEventType_Little {
		return "Little"
	}
	if e == NestsEventType_Elite {
		return "Elite"
	}
	if e == NestsEventType_Boss {
		return "Boss"
	}
	if e == NestsEventType_SecretRoom {
		return "SecretRoom"
	}
	if e == NestsEventType_Obstacle {
		return "Obstacle"
	}
	if e == NestsEventType_Trap {
		return "Trap"
	}
	if e == NestsEventType_Random {
		return "Random"
	}
	if e == NestsEventType_RedEnvelope {
		return "RedEnvelope"
	}
	if e == NestsEventType_TrapDebuff {
		return "TrapDebuff"
	}

	return ""

}

func GetNestsEventType(tp string) (NestsEventType, bool) {
	switch tp {
	case "None":
		return NestsEventType_None, true
	case "Award":
		return NestsEventType_Award, true
	case "Little":
		return NestsEventType_Little, true
	case "Elite":
		return NestsEventType_Elite, true
	case "Boss":
		return NestsEventType_Boss, true
	case "SecretRoom":
		return NestsEventType_SecretRoom, true
	case "Obstacle":
		return NestsEventType_Obstacle, true
	case "Trap":
		return NestsEventType_Trap, true
	case "Random":
		return NestsEventType_Random, true
	case "RedEnvelope":
		return NestsEventType_RedEnvelope, true
	case "TrapDebuff":
		return NestsEventType_TrapDebuff, true

	}
	return -1, false

}

func GetNestsEventType_ByAlias(tp string) (NestsEventType, bool) {
	switch tp {
	case "N":
		return NestsEventType_None, true
	case "A":
		return NestsEventType_Award, true
	case "B":
		return NestsEventType_Little, true
	case "C":
		return NestsEventType_Elite, true
	case "D":
		return NestsEventType_Boss, true
	case "E":
		return NestsEventType_SecretRoom, true
	case "Z":
		return NestsEventType_Obstacle, true
	case "T":
		return NestsEventType_Trap, true
	case "R":
		return NestsEventType_Random, true
	case "RE":
		return NestsEventType_RedEnvelope, true
	case "TD":
		return NestsEventType_TrapDebuff, true

	}
	return -1, false

}
