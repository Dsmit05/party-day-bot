package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const configFileName = "config.yml"

var (
	// buildVersion sets on compile time.
	buildVersion = ""
)

type Database struct {
	Use      bool   `yaml:"use"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Table    string `yaml:"table"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type ApiGRPCServer struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	TimeoutConnection int    `yaml:"timeoutConnection"`
}

type ApiHTTPServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Bot struct {
	Key       string `yaml:"key"`
	Debug     bool   `yaml:"debug"`
	SecretCMD string `yaml:"secretCMD"`
}

// Project - contains all parameters project information.
type Project struct {
	BuildVersion string
}

type Config struct {
	Bot           Bot           `yaml:"bot"`
	Database      Database      `yaml:"database"`
	ApiGRPCServer ApiGRPCServer `yaml:"apiGRPCServer"`
	ApiHTTPServer ApiHTTPServer `yaml:"apiHTTPServer"`
	Project
}

func NewConfig() (*Config, error) {
	var cfg = new(Config)
	if err := cfg.initFromFile(configFileName); err != nil {
		return nil, err
	}

	cfg.Project.BuildVersion = buildVersion

	return cfg, nil
}

// initFromFile init Config from yml file.
func (c *Config) initFromFile(filePath string) error {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("config file close: %v", err)
		}
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&c); err != nil {
		return err
	}

	return nil
}

// String representation Config settings.
func (c Config) String() string {
	return fmt.Sprintf(" BuildVersion: %+v\n Database: %+v\n Bot: %+v\n GRPCServer: %+v\n HttpServer: %+v\n",
		c.BuildVersion, c.Database, c.Bot, c.ApiGRPCServer, c.ApiGRPCServer)
}

func (c *Config) GetSecretCommand() string {
	return c.Bot.SecretCMD
}

func (c *Config) IsDBUse() bool {
	return c.Database.Use
}

func (c *Config) GetConnectDB() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Table,
	)
}

func (c *Config) GetApiGRPCServerAddr() string {
	return fmt.Sprintf("%v:%v", c.ApiGRPCServer.Host, c.ApiGRPCServer.Port)
}

func (c *Config) GetApiGRPCServerTimeout() time.Duration {
	timeout := time.Duration(c.ApiGRPCServer.TimeoutConnection) * time.Second
	return timeout
}

func (c *Config) GetApiHTTPServerAddr() string {
	return fmt.Sprintf("%v:%v", c.ApiHTTPServer.Host, c.ApiHTTPServer.Port)
}
