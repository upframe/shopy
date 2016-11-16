package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Development    bool
	Key            string
	Domain         string
	Port           int
	Scheme         string
	Assets         string
	InviteOnly     bool
	DefaultInvites int
	Database       struct {
		User     string
		Password string
		Host     string
		Port     string
		Name     string
	}
	SMTP struct {
		User     string
		Password string
		Host     string
		Port     string
	}
	PayPal struct {
		Client string
		Secret string
	}
}

func configFile(path string) (*config, error) {
	file := &config{}

	configFile, err := os.Open("config.json")
	if err != nil {
		return file, err
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&file)
	return file, nil
}
