package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

/*
Config is used in user system.
*/
type Config struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	SessionID    string `yaml:"session_id"`
	Pin          string `yaml:"pin"`
	PinToken     string `yaml:"pin_token"`
	PrivateKey   string `yaml:"private_key"`
}

var cfg *Config

func LoadConfig() (*Config, error) {
	cfg = new(Config)
	bytes, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(bytes), cfg)
	if err != nil {
		panic(err)
	}
	return cfg, nil
}

func GetConfig() *Config {
	var err error
	if cfg == nil {
		if cfg, err = LoadConfig(); err != nil {
			log.Panicf("Failed to load config: %s\n", err)
		}
	}
	return cfg
}
