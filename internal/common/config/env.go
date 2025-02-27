package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func MustLoadEnv(filePath string) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		panic(fmt.Errorf("cannot stat file %s: %w", filePath, err))
	}

	err = godotenv.Load(filePath)
	if err != nil {
		panic(fmt.Errorf("error loading .env file: %w", err))
	}
}
