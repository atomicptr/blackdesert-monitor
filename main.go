package main

import (
	"log"
	"time"

	"github.com/atomicptr/blackdesert-monitor/system"
)

const BlackDesertProcessName64bit = "BlackDesert64.exe"

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	for {
		processId, err := system.FindProcessByName(BlackDesertProcessName64bit)
		if err != nil {
			return err
		}

		conns, _ := system.GetConnectionsWithPid(processId)

		log.Printf("BDO PID: %d, Connections: %d, RUNNING? %v\n", processId, len(conns), len(conns) > 0)

		time.Sleep(30 * time.Second)
	}
}
