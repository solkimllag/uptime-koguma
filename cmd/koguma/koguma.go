package koguma

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"uptime-koguma/pkg/systemchecks"
)

var conf Config
var status string

func loadConfFile(filePath *string) error {
	rawConf, err := os.ReadFile(*filePath)
	if err != nil {
		return err
	}
	json.Unmarshal(rawConf, &conf)
	return nil
}

func Koguma() {
	confFile := flag.String("f", "", "Specify config file path.")
	flag.Parse()
	if *confFile == "" {
		*confFile = os.Getenv("ENV_KOGUMA_CONF")
	}
	err := loadConfFile(confFile)
	if err != nil {
		log.Fatal(err)
	}

	for {
		status = "up"
		if conf.CPUThreshold > 0 {
			cpuLoad, err := systemchecks.GetLoadAvarage(conf.CPULoadAveragaeType)
			if err != nil {
				log.Print(err)
			}
			if cpuLoad >= conf.CPUThreshold {
				status = "down"
			}
		}
		if conf.MemoryThreshold > 0 {
			freeMem, err := systemchecks.GetFreeMem()
			if err != nil {
				log.Print(err)
			}
			if freeMem <= uint64(conf.MemoryThreshold) {
				status = "down"
			}
		}
		for _, d := range conf.Disks {
			freeDisk, err := systemchecks.GetFreeSpace(d.Path)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s %d\n", d.Path, freeDisk)
			if freeDisk <= d.Threshold {
				status = "down"
			}
		}
		err := SendHeartbeat()
		if err != nil {
			log.Print(err)
		}
		time.Sleep(time.Duration(conf.HeartbeatInterval) * time.Second)
	}
}

// Sends a heartbeat with the status 'up' or 'down' to uptime-kuma server.
func SendHeartbeat() error {
	pushUrl := conf.PushURL
	if status == "down" {
		pushUrl = strings.Replace(pushUrl, "up", "down", 1)
	}
	resp, err := http.Get(pushUrl)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return errors.New(string(b))
	}
	return nil
}
