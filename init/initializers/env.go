package initializers

import (
	"os"
)

const defaultEnv = "development"

func Env() string {
	goEnv := os.Getenv("GO_ENVIRONMENT")
	if goEnv == "" {
		return defaultEnv // default environment
	}
	return goEnv
}