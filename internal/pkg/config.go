package pkg

import (
	"io"
	"log/slog"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   Server `yaml:"server"`
	Database DB     `yaml:"database"`
	AWS      Aws    `yaml:"aws"`
}

type Server struct {
	PORT int `yaml:"port"`
}

type Aws struct {
	Region string `yaml:"region"`
}

type DB struct {
	DATABASE_URL string `yaml:"database_url"`
}

func (c *Config) LoadConfig(file io.Reader) error {
	data, err := io.ReadAll(file)
	if err != nil {
		slog.Error("failed to read the provided yaml file")
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		slog.Error("failed to umarshal file data", "ERR", err)
		return err
	}

	return nil
}
