package command

import "cook-robot-middle-platform-go/instruction"

const (
	COOK         = "cook"         // multiple
	WASH         = "wash"         // multiple
	RESET        = "reset"        // multiple
	DOOR_UNLOCK  = "door_unlock"  // single
	DISH_OUT     = "dish_out"     // multiple
	RESUME       = "resume"       // single
	PAUSE_TO_ADD = "pause_to_add" // single
)

const (
	MULTIPLE = "multiple"
	SINGLE   = "single"
)

type Command struct {
	CommandType  string                      `json:"commandType"`
	CommandName  string                      `json:"commandName"`
	Instructions []instruction.Instructioner `json:"instructions"`
}
