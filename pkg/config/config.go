package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	Domain              string `yaml:"domain"`
	IngressURLTemplate  string `yaml:"ingressURLTemplate"`
	IngressNameTemplate string `yaml:"ingressNameTemplate"`
}

//ReadConfig function that reads the yaml file
func ReadConfig(filePath string) Configuration {
	var config Configuration
	// Read YML
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}

	// Unmarshall
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Panic(err)
	}

	return config
}

//WriteConfig function that can write to the yaml file
func WriteConfig(config Configuration, path string) error {
	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}

	return nil
}
