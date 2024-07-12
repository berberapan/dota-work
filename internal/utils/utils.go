package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

func ConvertIntToBinary(number int) string {
	return fmt.Sprintf("%b", number)
}

func ConvertDateToUnix(date string) int {
	layout := "2006-01-02"
	utcTime, err := time.Parse(layout, date)
	if err != nil {
		log.Fatalf("Unable to convert date to time\nError: %s", err)
	}
	return int(utcTime.Unix())
}

func ConvertSecondsToFormattedString(duration int) string {
	minutes := duration / 60
	seconds := duration % 60
	return fmt.Sprintf("%dm%ds", minutes, seconds)
}

func GetEnvVariable(key string) string {
	path, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path.\nError message: %s", err)
	}
	dir := filepath.Dir(path)
	envPath := filepath.Join(dir, "../", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Unable to load .env file.\nError message: %s", err)
	}
	return os.Getenv(key)
}
