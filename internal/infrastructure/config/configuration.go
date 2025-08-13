package config

import (
	"fmt"
	"strings"

	"sync"

	"github.com/Ifelsik/url-shortener/internal/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type Server struct {
	Host string
	Port string
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Logger struct {
	Level      string
	Formatter  string
	ShowCaller bool
}

type Config struct {
	Server   Server
	Database Database
	Logger   Logger
}

type configServer struct {
	config *Config
	koanf  *koanf.Koanf
	mu     *sync.RWMutex
}

func NewConfigServer() *configServer {
	return &configServer{
		koanf:  koanf.New("."),
		config: new(Config),
		mu:     &sync.RWMutex{},
	}
}

func (c *configServer) Load(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("read env file: %w", err)
	}

	err = c.koanf.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		return fmt.Errorf("read yaml file: %w", err)
	}

	err = c.koanf.Load(
		env.Provider(
			"", ".", func(s string) string {
				return strings.ReplaceAll(strings.ToLower(s), "_", ".")
			}),
		nil)
	if err != nil {
		return fmt.Errorf("read env variables: %w", err)
	}

	err = c.koanf.Unmarshal("", c.config)
	if err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	return nil
}

func (c *configServer) GetServerConfig() Server {
	return c.config.Server
}

func (c *configServer) GetDatabaseConfig() Database {
	return c.config.Database
}

func (c *configServer) GetLoggerConfig() logger.LoggerConfig {	
	var level uint8
	switch c.config.Logger.Level {
	case "error":
		level = logger.LevelError
	case "warning":
		level = logger.LevelWarning
	case "info":
		level = logger.LevelInfo
	case "debug":
		level = logger.LevelDebug
	default:
		level = logger.LevelInfo
	}
	
	var formatter uint8
	switch c.config.Logger.Formatter {
	case "text":
		formatter = logger.TextFormatter
	case "json":
		formatter = logger.JSONFormatter
	default:
		formatter = logger.TextFormatter
	}
	
	logConf := logger.LoggerConfig{
		Level:      level,
		Formatter:  formatter,
		ShowCaller: c.config.Logger.ShowCaller,
	}

	return  logConf
}
