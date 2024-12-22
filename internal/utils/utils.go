package utils

import (
	"fmt"
	"log"
	"os"
)

func MakeLogger(filename string) (*log.Logger, error) {
	logFile, err := os.OpenFile(fmt.Sprintf("log/%s.log", filename), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	logger := log.New(logFile, "", log.LstdFlags)
	return logger, nil
}
