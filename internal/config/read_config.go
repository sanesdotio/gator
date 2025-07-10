package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func Read() *Config{
	configFilePath, err := getConfigFilePath(configFileName)
	if err != nil {
		fmt.Println("Error getting config file path:", err)
		return &Config{}
	}
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return &Config{}
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("Error unmarshalling config file:", err)
		return &Config{}
	}

	return &config
}