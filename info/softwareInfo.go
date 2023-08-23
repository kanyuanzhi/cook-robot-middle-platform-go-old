package info

import (
	"cook-robot-middle-platform-go/utils"
	"time"
)

var Software = &SoftwareInfo{}

type SoftwareInfo struct {
	Name         string    `json:"name" yaml:"name" mapstructure:"name"`
	Version      string    `json:"version" yaml:"version" mapstructure:"version"`
	MachineModel string    `json:"machineModel" yaml:"machineModel" mapstructure:"machineModel"`
	UpdateTime   time.Time `json:"updateTime" yaml:"updateTime" mapstructure:"updateTime"`
}

func (info *SoftwareInfo) Reload() {
	utils.Reload("softwareInfo", Software)
}

func init() {
	Software.Reload()
}
