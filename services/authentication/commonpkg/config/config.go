package config

import (
	"path/filepath"

	"github.com/joho/godotenv"
)

func Loadconfig(path string) (err error) {
	err = godotenv.Load(filepath.Join(path, ".env"))
	if err != nil {
		return err
	}
	return nil
}
