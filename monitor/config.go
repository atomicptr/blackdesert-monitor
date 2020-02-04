package monitor

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ProcessName  string        `yaml:"processName"`
	PollInterval time.Duration `yaml:"pollInterval"`
	Telegram     struct {
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
