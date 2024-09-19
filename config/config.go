package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type RepositoryConfig struct {
	Path   string `json:"path"`
	Branch string `json:"branch"`
}

type Config struct {
	Port           int              `json:"port"`
	Repository     RepositoryConfig `json:"repository"`
	RestartCommand string           `json:"restart"`
	Secret         string           `json:"secret"`
}

func LoadConfig() Config {
	configFile := "./config/config.json"

	f, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = json.Unmarshal(d, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
