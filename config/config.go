package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LoginIp    string `yaml:"login_ip"`
	LoginPort  int    `yaml:"login_port"`
	VersionMin int    `yaml:"version_min"`
	VersionMax int    `yaml:"version_max"`
	GameIp     string `yaml:"game_ip"`
	GamePort   int    `yaml:"game_port"`
}

var ConfigInstance = &Config{}

func (c *Config) Load() {
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		log.Fatal(err)
	}
}
