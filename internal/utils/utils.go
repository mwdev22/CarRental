package utils

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func MakeLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(fmt.Sprintf("log/%s.log", filename), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return log.Default()
	}
	logger := log.New(logFile, "", log.LstdFlags)
	return logger
}

func GenerateUniqueString(base string) string {
	return fmt.Sprintf("%s_%d", base, time.Now().UnixNano()+int64(rand.Intn(1000)))
}
