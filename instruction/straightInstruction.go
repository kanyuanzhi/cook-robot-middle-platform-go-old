package instruction

type function string

const (
	RESET  = function("reset")
	LOCATE = function("locate")
	START  = function("start")
	STOP   = function("stop")
)

type AxisInstruction struct {
	Instruction    `mapstructure:",squash"`
	Function       function `json:"function"`
	Axis           string   `json:"AXIS"`
	TargetPosition uint32   `json:"target_position"  mapstructure:"target_position"`
	Speed          uint32   `json:"speed"`
}

func NewAxisInstruction(function function, axis string, targetPosition uint32, speed uint32) *AxisInstruction {
	return &AxisInstruction{
		Instruction:    NewInstruction(AXIS),
		Function:       function,
		Axis:           axis,
		TargetPosition: targetPosition,
		Speed:          speed,
	}
}

type RotateInstruction struct {
	Instruction      `mapstructure:",squash"`
	Function         function `json:"function"`
	Mode             uint32   `json:"mode"`
	Speed            uint32   `json:"speed"`
	RotationalAmount uint32   `json:"rotational_amount" mapstructure:"rotational_amount"`
}

func NewRotateInstruction(name string, function function, mode uint32, speed uint32, rotationalAmount uint32) *RotateInstruction {
	return &RotateInstruction{
		Instruction: Instruction{
			InstructionType: ROTATE,
			InstructionName: name,
		},
		Function:         function,
		Mode:             mode,
		Speed:            speed,
		RotationalAmount: rotationalAmount,
	}
}

type PumpInstruction struct {
	Instruction `mapstructure:",squash"`
	PumpNumber  uint32 `json:"pump_number" mapstructure:"pump_number"`
	Duration    uint32 `json:"duration"`
}

func NewPumpInstruction(pumpNumber uint32, duration uint32) *PumpInstruction {
	return &PumpInstruction{
		Instruction: NewInstruction(PUMP),
		PumpNumber:  pumpNumber,
		Duration:    duration,
	}
}
