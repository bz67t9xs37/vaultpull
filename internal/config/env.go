package config

import "os"

// lookupEnv is a thin wrapper around os.Getenv to allow test overriding.
var lookupEnv = func(key string) string {
	return os.Getenv(key)
}
