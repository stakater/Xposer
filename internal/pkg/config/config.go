package config

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	Domain              string `yaml:"domain"`
	IngressURLTemplate  string `yaml:"ingressURLTemplate"`
	IngressURLPath      string `yaml:"ingressURLPath"`
	IngressNameTemplate string `yaml:"ingressNameTemplate"`
	TLS                 bool   `yaml:"tls"`
	ExposeServiceUrl    string `yaml:"exposeServiceURL"`
}

//ReadConfig function that reads the yaml file
func ReadConfig(filePath string) (Configuration, error) {
	var config Configuration
	// Read YML
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("Error reading config: %v", err)
		return config, err
	}

	// Unmarshall
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		logrus.Errorf("Error unmarshalling config: %v", err)
		return config, err
	}

	return config, nil
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

func GetControllerConfig() Configuration {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		configFilePath = "configs/config.yaml"
	}

	configuration, err := ReadConfig(configFilePath)
	if err != nil {
		logrus.Errorf("Can not read configuration file with the following error: %v", err)
	}
	return configuration
}
