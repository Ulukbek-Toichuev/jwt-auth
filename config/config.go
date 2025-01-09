package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const path = "../../../config.yml"

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		DriverName string `yaml:"driver_name"`
		DataSource string `yaml:"data_source"`
	} `yaml:"database"`
	Jwt struct {
		SecretKey   string `yaml:"secret_key"`
		TokenExpiry int    `yaml:"token_expiry"`
	} `yaml:"jwt"`
}

func NewConfig() *Config {
	var config Config
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}
	defer file.Close()

	decode := yaml.NewDecoder(file)
	if err := decode.Decode(&config); err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}

	return &config
}

func (cg *Config) GetPort() string {
	return cg.Server.Port
}

func (cg *Config) GetDriverName() string {
	return cg.Database.DriverName
}

func (cg *Config) GetDataSource() string {
	return cg.Database.DataSource
}

func (cg *Config) GetSecretKey() string {
	return cg.Jwt.SecretKey
}

func (cg *Config) GetTokenExpiry() int {
	return cg.Jwt.TokenExpiry
}
