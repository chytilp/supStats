package common

import (
	"io"
	"os"

	toml "github.com/pelletier/go-toml"
)

var config *Config

type Config struct {
	BackendUrl        string
	FrontendUrl       string
	MobileUrl         string
	DataFolder        string
	CategoryBaseUrl   string
	TechnologyBaseUrl string
	HeaderHost        string
	HeaderReferer     string
	HeaderOrigin      string
	DbPath            string
}

func GetConfig() *Config {
	if config == nil {
		config = read()
	}
	return config
}

func read() *Config {
	file, err := os.Open("configProd.toml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config Config

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}
	return &config
}
