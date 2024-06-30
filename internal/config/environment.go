package config

import "os"

func IsLocal() bool {
	env := os.Getenv("ENVIRONMENT_NAME")
	return env == "local" || env == ""
}
