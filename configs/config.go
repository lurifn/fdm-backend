package configs

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Configuration struct {
	Email struct {
		From struct {
			Email    string `yaml:"email"`
			Password string `yaml:"password"`
			SMTP     string `yaml:"smtp"`
			Port     string `yaml:"port"`
		} `yaml:"from"`
		To []string `yaml:"to"`
	} `yaml:"email"`
}

func Load() (*Configuration, error) {
	var c Configuration
	f, err := os.Open("configs/config.yml")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Println("error closing config yml: ", err.Error())
		}
	}()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}