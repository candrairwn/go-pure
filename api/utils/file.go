package utils

import (
	"os"

	"go.uber.org/zap"
)

func ReadFile(path string, log *zap.SugaredLogger) (string, error) {

	// get current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Error(err)
	}

	// get file path
	path = dir + path

	value, err := os.ReadFile(path)
	if err != nil {
		log.Error(err)
	}
	return string(value), nil
}

func ReadFileReturnByte(path string, log *zap.SugaredLogger) ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Error(err)
	}

	// get file path
	path = dir + path

	value, err := os.ReadFile(path)
	if err != nil {
		log.Error(err)
	}
	return value, nil
}
