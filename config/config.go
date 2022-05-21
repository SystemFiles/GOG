package config

import (
	"bytes"
	"errors"
	"io"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
	"sykesdev.ca/gog/internal/common"
	"sykesdev.ca/gog/internal/logging"
)

var defaults = `
---
# GOG Configuration Settings

logging:
  level: "INFO"
application:
  tag_prefix: "v"
`

var (
	once sync.Once
	instance Configuration
)

type Configuration struct {
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
	
	Application struct {
		TagPrefix string `yaml:"tag_prefix"`
	} `yaml:"application"`
}

func AppConfig() *Configuration {
	once.Do(func ()  {
		instance = Configuration{}

		err := instance.load()
		if err != nil {
			panic("failed to load app configuration instance\n" + err.Error())
		}
	})

	return &instance
}

func (c *Configuration) load() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appConfigPath := configDir + "/gog/config.yml"

	var f *os.File
	if !common.PathExists(appConfigPath) {
		if err := os.MkdirAll(configDir + "/gog/", 0755); err != nil { return err }
		f, err = os.Create(appConfigPath)
		if err != nil {
			return err
		}
		
		// pre-load with defaults
		_, err = f.Write([]byte(defaults))
		if err != nil {
			return err
		}
	} else {
		f, err = os.Open(appConfigPath)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		return err
	}

	err = yaml.Unmarshal(buf.Bytes(), &c)
	return err
}

func (c *Configuration) TagPrefix() string {
	return c.Application.TagPrefix
}

func (c *Configuration) SetTagPrefix(prefix string) {
	c.Application.TagPrefix = prefix
}

func (c *Configuration) LogLevel() string {
	return c.Logging.Level
}

func (c *Configuration) SetLogLevel(level string) error {
	if !common.StringInSlice(logging.SeverityLevels, level) {
		return errors.New("error setting log severity level. invalid severity level provided")
	}

	c.Logging.Level = level
	return nil
}