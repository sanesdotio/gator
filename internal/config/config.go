package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(userName string) error {
	if userName == "" {
		return nil
	}
	cfg.CurrentUserName = userName
	
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	write(jsonData)
	return nil
}

func write(jsonData []byte) error {
	filePath, err := getConfigFilePath(configFileName);
	if err != nil {
		return err
	}
	// Write the JSON data to the file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
			return err
		}
	return nil
}

func getConfigFilePath(configFileName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory")
		return "", err
	}
	return fmt.Sprintf("%s/%s", homeDir, configFileName), nil
}