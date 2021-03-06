package system

import (
	"github.com/mitchellh/go-ps"
	"github.com/pkg/errors"
	"os"
)

func FindProcessByName(processName string) (int, error) {
	processes, err := ps.Processes()
	if err != nil {
		return -1, errors.Wrap(err, "could not read processes")
	}

	for _, process := range processes {
		if process.Executable() == processName {
			return process.Pid(), nil
		}
	}

	return -1, errors.Errorf("could not find system \"%s\"", processName)
}

func Kill(processId int) error {
	process, err := os.FindProcess(processId)
	if err != nil {
		return err
	}

	return process.Kill()
}
