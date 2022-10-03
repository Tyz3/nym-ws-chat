package config

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

//go:embed resource/config.yaml
var BinConfig []byte

type Contact struct {
	Address string `yaml:"address"`
	Alias   string `yaml:"alias"`
}

type Config struct {
	configName string

	Client struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"Client"`
	Contacts []Contact `yaml:"Contacts"`
}

func NewConfig(configName string) *Config {
	cfg := &Config{configName: configName}
	data, err := os.ReadFile(configName)
	if err != nil {
		panic(err)
	}

	// Парсим конфиг
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func (c *Config) Save() {
	file, err := os.OpenFile(c.configName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	data, err := yaml.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = file.Write(data)
	if err != nil {
		fmt.Println(err)
	}
}


