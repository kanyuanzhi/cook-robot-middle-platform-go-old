package instruction

type Instructioner interface {
	CheckType() InstructionType
	SetName(name string)
}

type InstructionType string

const (
	INGREDIENT       = InstructionType("ingredient")
	SEASONING        = InstructionType("seasoning")
	WATER            = InstructionType("water")
	OIL              = InstructionType("oil")
	STIR_FRY         = InstructionType("stir_fry")
	HEAT             = InstructionType("heat")
	DISH_OUT         = InstructionType("dish_out")
	SHAKE            = InstructionType("shake")
	LAMPBLACK_PURIFY = InstructionType("lampblack_purify")
	DOOR_UNLOCK      = InstructionType("door_unlock")
	RESET_ALL        = InstructionType("reset_all")
	RESET_XYT        = InstructionType("reset_xyt")
	RESET_RT         = InstructionType("reset_rt")
	PREPARE          = InstructionType("prepare")
	DELAY            = InstructionType("delay")
	RESUME           = InstructionType("resume")
	PAUSE_TO_ADD     = InstructionType("pause_to_add")

	AXIS   = InstructionType("axis")
	ROTATE = InstructionType("rotate")
	PUMP   = InstructionType("pump")
)

type Instruction struct {
	InstructionType InstructionType `json:"instructionType" mapstructure:"instructionType"`
	InstructionName string          `json:"instructionName" mapstructure:"instructionName"`
}

func (i Instruction) CheckType() InstructionType {
	return i.InstructionType
}

func (i Instruction) SetName(name string) {
	i.InstructionName = name
}

func NewInstruction(instructionType InstructionType) Instruction {
	return Instruction{InstructionType: instructionType}
}

type IngredientInstruction struct {
	Instruction `mapstructure:",squash"`
	SlotNumber  string `json:"slotNumber" mapstructure:"slotNumber"`
}

func NewIngredientInstruction(slotNumber string) *IngredientInstruction {
	return &IngredientInstruction{
		Instruction: NewInstruction(INGREDIENT),
		SlotNumber:  slotNumber,
	}
}

type SeasoningInstruction struct {
	Instruction     `mapstructure:",squash"`
	PumpToWeightMap map[string]uint32 `json:"pumpToWeightMap" mapstructure:"pumpToWeightMap"` // 泵号:重量
}

func NewSeasoningInstruction(name string, pumpToWeightMap map[string]uint32) *SeasoningInstruction {
	return &SeasoningInstruction{
		Instruction: Instruction{
			InstructionType: SEASONING,
			InstructionName: name,
		},
		PumpToWeightMap: pumpToWeightMap,
	}
}

type WaterInstruction struct {
	Instruction `mapstructure:",squash"`
	PumpNumber  uint32 `json:"pumpNumber" mapstructure:"pumpNumber"`
	Weight      uint32 `json:"weight"`
}

func NewWaterInstruction(pumpNumber uint32, weight uint32) *WaterInstruction {
	return &WaterInstruction{
		Instruction: NewInstruction(WATER),
		PumpNumber:  pumpNumber,
		Weight:      weight,
	}
}

type OilInstruction struct {
	Instruction `mapstructure:",squash"`
	PumpNumber  uint32 `json:"pumpNumber" mapstructure:"pumpNumber"`
	Weight      uint32 `json:"weight"`
}

func NewOilInstruction(pumpNumber uint32, weight uint32) *OilInstruction {
	return &OilInstruction{
		Instruction: NewInstruction(OIL),
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

type HeatInstruction struct {
	Instruction       `mapstructure:",squash"`
	Temperature       float64 `json:"temperature"`
	TargetTemperature float64 `json:"targetTemperature" mapstructure:"targetTemperature"`
	Duration          uint32  `json:"duration"`
	JudgeType         uint    `json:"judgeType" mapstructure:"judgeType"`
}

func NewHeatInstruction(temperature float64, targetTemperature float64, duration uint32, judgeType uint) *HeatInstruction {
	return &HeatInstruction{
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
	NO_JUDGE
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

type ResetAllInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewResetAllInstruction() *ResetAllInstruction {
	return &ResetAllInstruction{
		Instruction: NewInstruction(RESET_ALL),
	}
}

type ResetXYTInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewResetXYTInstruction() *ResetXYTInstruction {
	return &ResetXYTInstruction{
		Instruction: NewInstruction(RESET_XYT),
	}
}

type ResetRTInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewResetRTInstruction(name string) *ResetRTInstruction {
	return &ResetRTInstruction{
		Instruction: Instruction{
			InstructionType: RESET_RT,
			InstructionName: name,
		},
	}
}

type PrepareInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewPrepareInstruction() *PrepareInstruction {
	return &PrepareInstruction{
		Instruction: NewInstruction(PREPARE),
	}
}

type DelayInstruction struct {
	Instruction `mapstructure:",squash"`
	Duration    uint32 `json:"duration"`
}

func NewDelayInstruction(duration uint32) *DelayInstruction {
	return &DelayInstruction{
		Instruction: NewInstruction(DELAY),
		Duration:    duration,
	}
}

type ResumeInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewResumeInstruction() *ResumeInstruction {
	return &ResumeInstruction{
		Instruction: NewInstruction(RESUME),
	}
}

type PauseToAddInstruction struct {
	Instruction `mapstructure:",squash"`
}

func NewPauseToAddInstruction() *PauseToAddInstruction {
	return &PauseToAddInstruction{
		Instruction: NewInstruction(PAUSE_TO_ADD),
	}
}

var InstructionTypeToStruct = map[InstructionType]Instructioner{
	INGREDIENT:       IngredientInstruction{},
	SEASONING:        SeasoningInstruction{},
	WATER:            WaterInstruction{},
	OIL:              OilInstruction{},
	STIR_FRY:         StirFryInstruction{},
	HEAT:             HeatInstruction{},
	DISH_OUT:         DishOutInstruction{},
	SHAKE:            ShakeInstruction{},
	LAMPBLACK_PURIFY: LampblackPurifyInstruction{},
	DOOR_UNLOCK:      DoorUnlockInstruction{},
	RESET_ALL:        ResetAllInstruction{},
	RESET_XYT:        ResetXYTInstruction{},
	RESET_RT:         ResetRTInstruction{},
	PREPARE:          PrepareInstruction{},
	DELAY:            DelayInstruction{},
	RESUME:           ResumeInstruction{},
	PAUSE_TO_ADD:     PauseToAddInstruction{},

	AXIS:   AxisInstruction{},
	ROTATE: RotateInstruction{},
	PUMP:   PumpInstruction{},
}
