package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func Parse(path string, cfg interface{}) error {
	bs, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bs, cfg)
	if err != nil {
		return err
	}
	return nil
}
