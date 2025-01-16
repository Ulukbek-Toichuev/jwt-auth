package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const envName string = "CONFIG_PATH"

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		DriverName     string `yaml:"driver_name"`
		DataSource     string `yaml:"data_source"`
		EnvName        string `yaml:"env_name"`
		MigrationsPath string `yaml:"migrations_path"`
	} `yaml:"database"`
	Jwt struct {
		SecretKey   string `yaml:"secret_key"`
		TokenExpiry int    `yaml:"token_expiry"`
	} `yaml:"jwt"`
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	path := getPath(envName)
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
	config.Database.MigrationsPath = getPath(config.Database.EnvName)
	config.Jwt.SecretKey = getPath("SECRET_KEY")
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

func (cg *Config) GetMigrationsPath() string {
	return cg.Database.MigrationsPath
}

func (cg *Config) GetSecretKey() string {
	return cg.Jwt.SecretKey
}

func (cg *Config) GetTokenExpiry() int {
	return cg.Jwt.TokenExpiry
}

func getPath(envName string) string {
	path := os.Getenv(envName)
	if path == "" {
		log.Fatalf("Cant find env variable: %s please check name", envName)
		os.Exit(1)
	}

	return path
}
