// Package config provides the means to load and decode configuration files.
package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

func LoadFromYAML(path string, target interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(target)
	if err != nil {
		return fmt.Errorf("decode config: %w", err)
	}
	err = validator.New().Struct(target)
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	return nil
}
