package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No.env file found")
	}
}

func GetEnv(key string) string {
	val, isExist := os.LookupEnv(key)
	if !isExist {
		return ""
	}
	return val
}
