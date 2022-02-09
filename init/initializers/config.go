package initializers

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
)

var configs map[string][]byte

func loadConfigFromReader(reader io.Reader, ptr interface{}) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read the content from config reader. Err: %w", err)
	}

	if err := yaml.Unmarshal(content, ptr); err != nil {
		return fmt.Errorf("failed to unmarshal content from config reader. Err: %w", err)
	}
	return nil
}

func loadConfigFromFile(filePath string, ptr interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file at path %v. Err: %w", filePath, err)
	}

	if err := loadConfigFromReader(file, ptr); err != nil {
		return fmt.Errorf("failed loading the content of file %v. Err: %w", filePath, err)
	}
	return nil
}

func loadConfigForEnviroment(pointer interface{}) error {
	path, err := scopeConfigPath()
	if err != nil {
		return fmt.Errorf("failed find the env path: %s. Error: %w", path, err)
	}
	if err := loadConfigFromFile(path, pointer); err != nil {
		return fmt.Errorf("failed to load the config file for env path: %s. Error: %w", path, err)
	}
	return nil
}

func LoadConfigSection(section string, pointer interface{}) error {
	bytes, found := configs[section]
	if !found {
		return fmt.Errorf(`config section "%s" not found`, section)
	}
	if err := yaml.Unmarshal(bytes, pointer); err != nil {
		return fmt.Errorf("failed to load the config section. Err: %v", err)
	}
	return nil
}

func ConfigInitializer() {
	data := make(map[string]interface{})
	if err := loadConfigForEnviroment(data); err != nil {
		panic(err)
	}
	configs = make(map[string][]byte)
	for section, config := range data {
		bytes, err := yaml.Marshal(config)
		if err != nil {
			panic(err)
		}
		configs[section] = bytes
	}
}