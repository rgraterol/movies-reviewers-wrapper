package initializers

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

const (
	confDefaultDir = "config/resources"
	appDefaultPath = "/init"
	ymlExtension = ".yml"
	yamlExtension = ".yaml"
)

var (
	crrFSGetter fsGetter = osFSGetter
)

func appDir() string {
	for crrPath, _ := crrFSGetter.getwd(); crrPath != string(filepath.Separator); {
		if _, err := os.Stat(filepath.Join(crrPath, "go.mod")); err == nil {
			return crrPath
		}
		crrPath = filepath.Dir(crrPath)
	}
	return appDefaultPath
}

func configDir() (string, error) {
	confDir := confDefaultDir
	if !filepath.IsAbs(confDir) {
		confDir = filepath.Join(appDir(), confDir)
	}
	if info, err := os.Stat(confDir); err != nil || !info.IsDir() {
		return "", fmt.Errorf("configs dir not found at %v", confDir)
	}
	return confDir, nil
}

func scopeConfigPath() (string, error) {
	cfgPath, err := configDir()
	if err != nil {
		return "", err
	}

	cfgName := path.Join(cfgPath, Env())
	withYMLExt := fmt.Sprintf("%v%v", cfgName, ymlExtension)
	if _, err := os.Stat(withYMLExt); err == nil {
		return withYMLExt, nil
	}
	withYAMLExt := fmt.Sprintf("%v%v", cfgName, yamlExtension)
	if _, err := os.Stat(withYAMLExt); err == nil {
		return withYAMLExt, nil
	}
	return "", fmt.Errorf("missing scope (%v) config file", Env())
}