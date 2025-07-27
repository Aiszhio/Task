package utils

import (
	"os"

	errors_package "github.com/Aiszhio/Task/internal/errors"
)

func GetEnv(name string) (string, error) {
	dsn := os.Getenv(name)
	if dsn == "" {
		return "", errors_package.NoVariable
	}

	return dsn, nil
}
