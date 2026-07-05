package config

import (
	"io/ioutil"
	"os"

	"github.com/MESLEKDAA/tepegoz/pkg/model"
	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*model.Config, error) {

	file, err := os.Open(path)

	if err != nil {

		return nil, err
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)

	if err != nil {

		return nil, err
	}

	var cfg model.Config

	if err := yaml.Unmarshal(bytes, &cfg); err != nil {

		return nil, err
	}

	return &cfg, nil

}
