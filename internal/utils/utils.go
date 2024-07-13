package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func ConvertDateToUnix(date string) int {
	layout := "2006-01-02"
	utcTime, err := time.Parse(layout, date)
	if err != nil {
		log.Fatalf("Unable to convert date to time\nError: %s", err)
	}
	return int(utcTime.Unix())
}

func ConvertSecondsToFormattedString(duration float64) string {
	minutes := int(duration) / 60
	seconds := int(duration) % 60
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

func ConvertBinIntToCount(number, width int) int {
	bin := fmt.Sprintf("%0*b", width, number)
	return strings.Count(bin, "0")
}

func CalculateAverageFromIntSlice(numbers []int) float64 {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return float64(sum) / float64(len(numbers))
}

func CalculateMedianFromIntSlice(numbers []int) float64 {
	slices.Sort(numbers)
	length := len(numbers)
	if length%2 == 0 {
		return float64((numbers[length/2-1] + numbers[length/2]) / 2)
	}
	return float64(numbers[length/2])
}

func CalculateAverageFromFloatSlice(numbers []float64) float64 {
	sum := 0.0
	for _, n := range numbers {
		sum += n
	}
	return float64(sum) / float64(len(numbers))
}

func CalculateMedianFromFloatSlice(numbers []float64) float64 {
	slices.Sort(numbers)
	length := len(numbers)
	if length%2 == 0 {
		return float64((numbers[length/2-1] + numbers[length/2]) / 2)
	}
	return float64(numbers[length/2])
}

func CalculatePercentageFromIntSlice(numbers []int) float64 {
	count := 0
	sum := 0
	for _, n := range numbers {
		if n == -1 {
			continue
		}
		count++
		sum += n
	}
	return float64(sum) / float64(count)
}
