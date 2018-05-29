package main

import (
	"github.com/PROger4ever-Golang/Redis-serialization-benchmarks/utils"
	"github.com/jinzhu/configor"
)

type Address struct {
	Host string `yaml:"Host"`
	Port int    `yaml:"Port"`
}

type Config struct {
	CommonRedis Address `yaml:"CommonRedis" required:"true"`
	RejsonRedis Address `yaml:"RejsonRedis" required:"true"`
}

func LoadConfig(files ...string) (*Config, error) {
	existingFiles, err := utils.GetExistingFiles(files...)
	if err != nil {
		return nil, utils.WrapIfError(err, "config->LoadConfig()->GetExistingFiles()")
	}

	var config Config
	err = configor.New(&configor.Config{ENVPrefix: "-"}).Load(&config, existingFiles...)
	if err != nil {
		return &config, utils.WrapIfError(err, "config->LoadConfig()->Load()")
	}
	return &config, nil
}
