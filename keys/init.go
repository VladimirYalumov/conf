package keys

import (
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v3"
)

type SourceType string

const (
	Yaml = SourceType("yaml")
	Env  = SourceType("env")
)

type DataType struct {
	Configs map[string]ConfigType `yaml:"configs"`
}

type ConfigType struct {
	Type    string     `yaml:"type"`
	Default string     `yaml:"default"`
	Source  SourceType `yaml:"source"`
}

func Init(r io.Reader) (confKeys map[SourceType]*DataType, err error) {
	data := new(DataType)
	if err = yaml.NewDecoder(r).Decode(data); err != nil {
		return nil, fmt.Errorf("keys: cannot decode keys yaml: %w", err)
	}

	configData := make(map[SourceType]*DataType, 2)

	configData[Yaml] = new(DataType)
	configData[Yaml].Configs = make(map[string]ConfigType)

	configData[Env] = new(DataType)
	configData[Env].Configs = make(map[string]ConfigType)

	for name, conf := range data.Configs {
		if conf.Source == Yaml {
			configData[conf.Source].Configs[name] = conf
		}
		if conf.Source == Env {
			configData[conf.Source].Configs[name] = conf
		}
	}
	return configData, nil
}
