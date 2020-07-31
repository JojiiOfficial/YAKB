package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config represents the structure
// of the config file
type Config struct {
	BotToken        string
	AllowSelfKarma  bool
	AddKarma        []string
	RemoveKarma     []string
	DataFile        string
	KarmaTopCommand string
	AllowBotVoting  bool
}

const configfile = "data/config.json"

func initConfig() (*Config, error) {
	s, err := os.Stat(configfile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if (s != nil && s.Size() == 0) || os.IsNotExist(err) {
		if err := createNewConfig(); err != nil {
			return nil, err
		}

		fmt.Println("New config created")
		os.Exit(0)
		return nil, nil
	}

	confData, err := ioutil.ReadFile(configfile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(confData, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func getDefaultConfig() Config {
	return Config{
		BotToken:        "PLACE TELEGRAM TOKEN HERE",
		AllowSelfKarma:  false,
		DataFile:        "data/data.db",
		AddKarma:        []string{"+1", "thx", "ty", "thankyou", "thanks", "thanx"},
		RemoveKarma:     []string{"-1", "rtfm"},
		KarmaTopCommand: "/ktop",
	}
}

func createNewConfig() error {
	return getDefaultConfig().Save()
}

// Save the config file
func (config Config) Save() error {
	d, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configfile, d, 0600)
}
