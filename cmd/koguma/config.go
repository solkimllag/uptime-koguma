package koguma

import (
	"time"
)

type Config struct {
	PushURL             string        `json:"push_url"`
	HeartbeatInterval   time.Duration `json:"heartbeat_interval"`
	CPUThreshold        uint          `json:"cpu_threshold"`
	CPULoadAveragaeType uint          `json:"cpu_load_average_type"` //1,5,15 min
	MemoryThreshold     uint          `json:"memory_threshold"`
	Disks               []Disk        `json:"disks"`
}

type Disk struct {
	Path      string `json:"disk_path"`
	Threshold uint   `json:"threshold"`
}
