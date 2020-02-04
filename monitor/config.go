package monitor

import (
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ProcessName                     string        `yaml:"processName"`
	PollInterval                    time.Duration `yaml:"pollInterval"`
	UnavailabilityThreshold         int           `yaml:"unavailabilityThreshold"`
	CloseBlackDesertWhenUnavailable bool          `yaml:"closeBlackDesertWhenUnavailable"`
	Telegram                        struct {
		Token  string `yaml:"token"`
		UserId int    `yaml:"userId"`
	}
}

func ConfigFromFile(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return configFromBytes(data)
}

func configFromBytes(data []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (config *Config) Validate() error {
	if config.ProcessName == "" {
		return errors.New("Process name can't be empty!")
	}

	if config.Telegram.Token == "" {
		return errors.New("Telegram Bot Token can't be empty!")
	}

	return nil
}
