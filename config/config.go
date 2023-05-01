package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LoginIp        string `yaml:"login_ip"`
	LoginPort      int    `yaml:"login_port"`
	VersionMin     int    `yaml:"version_min"`
	VersionMax     int    `yaml:"version_max"`
	GameIp         string `yaml:"game_ip"`
	GamePort       int    `yaml:"game_port"`
	ServerName     string `yaml:"server_name"`
	PemFile        string `yaml:"pem_file"`
	Motd           string `yaml:"motd"`
	DbHost         string `yaml:"db_host"`
	DbUser         string `yaml:"db_user"`
	DbPass         string `yaml:"db_pass"`
	DbDatabase     string `yaml:"db_database"`
	DbPort         int    `yaml:"db_port"`
	EncryptionType string `yaml:"encryption_type"`
}

var config = &Config{}

func Load() {
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConfig() *Config {
	return config
}
