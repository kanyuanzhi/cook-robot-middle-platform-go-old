package command

import "cook-robot-middle-platform-go/instruction"

const (
	COMMAND_NAME_COOK         = "cook"         // multiple
	COMMAND_NAME_WASH         = "wash"         // multiple
	COMMAND_NAME_POUR         = "pour"         // multiple
	COMMAND_NAME_PREPARE      = "prepare"      // multiple
	COMMAND_NAME_DOOR_UNLOCK  = "door_unlock"  // single
	COMMAND_NAME_DISH_OUT     = "dish_out"     // multiple
	COMMAND_NAME_RESUME       = "resume"       // single
	COMMAND_NAME_PAUSE_TO_ADD = "pause_to_add" // single
	COMMAND_NAME_HEAT         = "heat"         // single
	COMMAND_NAME_WITHDRAW     = "withdraw"     // multiple
)

const (
	COMMAND_TYPE_MULTIPLE = "multiple" // 不可在其他命令执行过程中执行
	COMMAND_TYPE_SINGLE   = "single"   // 可在其他命令执行过程中执行
)

type Command struct {
	CommandType  string                      `json:"commandType"`
	CommandName  string                      `json:"commandName"`
	DishUuid     string                      `json:"dishUuid"` //如果是炒制命令，则会携带菜品的uuid
	Instructions []instruction.Instructioner `json:"instructions"`
}
