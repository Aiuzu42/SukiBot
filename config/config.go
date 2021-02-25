package config

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	filePath = "sukiConfig.json"
)

type configuration struct {
	LogLevel   string      `json:"logLevel"`
	Token      string      `json:"token"`
	Users      UsersConfig `json:"users"`
	CustomSays []CustomCMD `json:"customSays"`
	Triggers   []Trigger   `json:"triggers"`
}

type UsersConfig struct {
	Owners []string `json:"owners"`
	Admins []string `json:"admins"`
}

type CustomCMD struct {
	Name    string `json:"commandName"`
	Channel string `json:"channel"`
}

type Trigger struct {
	Trigger       string   `json:"trigger"`
	Channels      []string `json:"channels"`
	Response      string   `json:"response"`
	Image         string   `json:"image"`
	Color         int      `json:"color"`
	CaseSensitive bool     `json:"caseSensitive"`
	Cooldown      int64    `json:"cooldown"`
	CooldownMap   map[string]int64
}

func (t Trigger) AllowsChannel(ch string) bool {
	if len(t.Channels) == 0 {
		return true
	}
	for _, e := range t.Channels {
		if e == ch {
			return true
		}
	}
	return false
}

var Config configuration
var Loc *time.Location

//InitConfig should be only used to load config at the start of the program, it panics if the config cannot be loaded for any reason.
func InitConfig() error {
	var err error
	Config, err = loadConfig()
	if err != nil {
		return err
	}
	switch Config.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	Loc = time.FixedZone("UTC-6", -6*60*60)
	return nil
}

//ReloadConfig can be used to reload config at any point, if it fails to reload it keeps the old config and returns an error.
func ReloadConfig() error {
	localConfig, err := loadConfig()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	Config = localConfig
	return nil
}

func loadConfig() (configuration, error) {
	var localConfig configuration
	file, err := os.Open(filePath)
	if err != nil {
		return configuration{}, errors.New("Unable to load configuration file [" + err.Error() + "]")
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&localConfig)
	if err != nil {
		return configuration{}, errors.New("Error parsing configuration file [" + err.Error() + "]")
	}
	for i := range localConfig.Triggers {
		localConfig.Triggers[i].CooldownMap = make(map[string]int64)
	}
	return localConfig, nil
}
