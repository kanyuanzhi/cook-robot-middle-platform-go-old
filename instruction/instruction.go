package instruction

type Instructioner interface {
	CheckType() InstructionType
}

type InstructionType string

const (
	INGREDIENT       = InstructionType("ingredient")
	SEASONING        = InstructionType("seasoning")
	WATER            = InstructionType("water")
	STIR_FRY         = InstructionType("stir_fry")
	HEAT             = InstructionType("heat")
	DISH_OUT         = InstructionType("dish_out")
	SHAKE            = InstructionType("shake")
	LAMPBLACK_PURIFY = InstructionType("lampblack_purify")
	DOOR_UNLOCK      = InstructionType("door_unlock")

	AXIS   = InstructionType("axis")
	ROTATE = InstructionType("rotate")
	PUMP   = InstructionType("pump")
)

type Instruction struct {
	InstructionType InstructionType `json:"instruction_type" mapstructure:"instruction_type"`
}

func (i Instruction) CheckType() InstructionType {
	return i.InstructionType
}

func NewInstruction(instructionType InstructionType) Instruction {
	return Instruction{InstructionType: instructionType}
}

type IngredientInstruction struct {
	Instruction `mapstructure:",squash"`
	SlotNumber  uint32 `json:"slot_number" mapstructure:"slot_number"`
}

func NewIngredientInstruction(slotNumber uint32) *IngredientInstruction {
	return &IngredientInstruction{
		Instruction: NewInstruction(INGREDIENT),
		SlotNumber:  slotNumber,
	}
}

type SeasoningInstruction struct {
	Instruction     `mapstructure:",squash"`
	PumpToWeightMap map[uint32]uint32 `json:"pump_to_weight_map" mapstructure:"pump_to_weight_map"` // 泵号:重量
}

func NewSeasoningInstruction(pumpToWeightMap map[uint32]uint32) *SeasoningInstruction {
	return &SeasoningInstruction{
		Instruction:     NewInstruction(SEASONING),
		PumpToWeightMap: pumpToWeightMap,
	}
}

type WaterInstruction struct {
	Instruction `mapstructure:",squash"`
	PumpNumber  uint32 `json:"pump_number" mapstructure:"pump_number"`
	Weight      uint32 `json:"weight"`
}

func NewWaterInstruction(pumpNumber uint32, weight uint32) *WaterInstruction {
	return &WaterInstruction{
		Instruction: NewInstruction(WATER),
		PumpNumber:  pumpNumber,
		Weight:      weight,
	}
}

type StirFryInstruction struct {
	Instruction `mapstructure:",squash"`
	Gear        uint32 `json:"gear"`
	Duration    uint32 `json:"duration"`
}

func NewStirFryInstruction(gear uint32, duration uint32) *StirFryInstruction {
	return &StirFryInstruction{
		Instruction: NewInstruction(STIR_FRY),
		Gear:        gear,
		Duration:    duration,
	}
}

type HeatingInstruction struct {
	Instruction       `mapstructure:",squash"`
	Temperature       float64 `json:"temperature"`
	TargetTemperature float64 `json:"target_temperature" mapstructure:"target_temperature"`
	Duration          uint32  `json:"duration"`
	JudgeType         uint    `json:"judge_type" mapstructure:"judge_type"`
}

func NewHeatingInstruction(temperature float64, targetTemperature float64, duration uint32, judgeType uint) *HeatingInstruction {
	return &HeatingInstruction{
		Instruction:       NewInstruction(HEAT),
		Temperature:       temperature,
		TargetTemperature: targetTemperature,
		Duration:          duration,
		JudgeType:         judgeType,
	}
}

const (
	BOTTOM_TEMPERATURE_JUDGE_TYPE uint = iota + 1
	INFRARED_TEMPERATURE_JUDGE_TYPE
	DURATION_JUDGE_TYPE
)

type DishOutInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewDishOutInstruction() *DishOutInstruction {
	return &DishOutInstruction{
		Instruction: NewInstruction(DISH_OUT),
	}
}

type ShakeInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewShakeInstruction() *DishOutInstruction {
	return &DishOutInstruction{
		Instruction: NewInstruction(SHAKE),
	}
}

const (
	VENTING uint32 = iota + 1
	PURIFICATION
)

type LampblackPurifyInstruction struct {
	Instruction `mapstructure:",squash"`
	Mode        uint32 `json:"mode"`
}

func NewLampblackPurifyInstruction(mode uint32) *LampblackPurifyInstruction {
	return &LampblackPurifyInstruction{
		Instruction: NewInstruction(LAMPBLACK_PURIFY),
		Mode:        mode,
	}
}

type DoorUnlockInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewDoorUnlockInstruction() *DoorUnlockInstruction {
	return &DoorUnlockInstruction{
		Instruction: NewInstruction(DOOR_UNLOCK),
	}
}

var InstructionTypeToStruct = map[InstructionType]Instructioner{
	INGREDIENT:       IngredientInstruction{},
	SEASONING:        SeasoningInstruction{},
	WATER:            WaterInstruction{},
	STIR_FRY:         StirFryInstruction{},
	HEAT:             HeatingInstruction{},
	DISH_OUT:         DishOutInstruction{},
	SHAKE:            ShakeInstruction{},
	LAMPBLACK_PURIFY: LampblackPurifyInstruction{},
	DOOR_UNLOCK:      DoorUnlockInstruction{},

	AXIS:   AxisInstruction{},
	ROTATE: RotateInstruction{},
	PUMP:   PumpInstruction{},
}
