package systemchecks

import (
	"runtime"
	"syscall"
)

const SI_LOAD_SHIFT = 16
const MB_SHIFT = 10

// GetLoadAvarage expects a LoadAvarage type and returns the correspondig load
// average value.
func GetLoadAvarage(loadav uint) (uint, error) {
	var load float64
	var sinfo syscall.Sysinfo_t
	err := syscall.Sysinfo(&sinfo)
	if err != nil {
		return 0, err
	}
	switch loadav {
	case 1:
		load = float64(sinfo.Loads[0])
	case 5:
		load = float64(sinfo.Loads[1])
	default: // case 15 min
		load = float64(sinfo.Loads[2])
	}
	return uint((load / float64(1<<SI_LOAD_SHIFT)) * 100 / float64(runtime.NumCPU())), nil
}

// GetFreeMem returns the available memory as a percentage of total system memory.
func GetFreeMem() (uint64, error) {
	var sinfo syscall.Sysinfo_t
	err := syscall.Sysinfo(&sinfo)
	if err != nil {
		return 0, err
	}
	availableRam := sinfo.Freeram>>MB_SHIFT + sinfo.Bufferram>>MB_SHIFT
	return uint64(float64(availableRam) / float64(sinfo.Totalram>>MB_SHIFT) * 100), nil
}

// GetFreeDisk expects a path, such as '/dev/sdb2' and returns the amount
// of remianing free disk space as a percentage of the full capacity.
func GetFreeSpace(path string) (uint, error) {
	var fs syscall.Statfs_t
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return 0, err
	}

	size := float64(fs.Blocks>>MB_SHIFT) * float64(fs.Bsize)
	free := float64(fs.Bfree>>MB_SHIFT) * float64(fs.Bsize)
	return uint(free / size * 100), nil
}
