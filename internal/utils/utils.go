package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func DataStructToJson[T any](data T) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal given struct.\n%w", err)
	}
	return jsonData, nil
}

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
	if seconds < 10 {
		return fmt.Sprintf("%dm0%ds", minutes, seconds)
	}
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

var months = map[string]time.Month{
	"January":   time.January,
	"February":  time.February,
	"March":     time.March,
	"April":     time.April,
	"May":       time.May,
	"June":      time.June,
	"July":      time.July,
	"August":    time.August,
	"September": time.September,
	"October":   time.October,
	"November":  time.November,
	"December":  time.December,
}

func ConvertStringToMonth(month string) time.Month {
	return months[month]
}

var timeZones = map[string]string{
	// Liquipedia abbr for timezones https://liquipedia.net/commons/Module:Timezone/Data
	// Europe
	"CEST": "Europe/Berlin",
	"CET":  "Europe/Berlin",
	"BST":  "Europe/London",
	"GMT":  "Europe/London",
	"UTC":  "Etc/UTC",
	"MSK":  "Europe/Moscow",
	"EET":  "Europe/Helsinki",
	"EEST": "Europe/Helsinki",
	// Asia
	"GST": "Asia/Dubai",
	"AST": "Asia/Riyadh",
	"SGT": "Asia/Singapore",
	"CST": "Asia/Shanghai",
	// Americas
	"EST":  "US/Eastern",
	"EDT":  "US/Eastern",
	"PST":  "US/Pacific",
	"PDT":  "US/Pacific",
	"CT":   "US/Central",
	"CDT":  "US/Central",
	"PET":  "America/Lima",
	"BRT":  "America/Sao_Paulo",
	"BRST": "America/Sao_Paulo",
}

func ConvertTzAbbrToTzString(abbr string) string {
	return timeZones[abbr]
}

func LiquipediaUrlToTournamentString(url string) (string, error) {
	cleanedUrl := strings.TrimSpace(url)
	_, tournament, ok := strings.Cut(cleanedUrl, "/dota2/")
	if !ok {
		return "", fmt.Errorf("Unable to use given URL.\nError: URL Format")
	}
	cleanedTournament := strings.ReplaceAll(tournament, "_", " ")
	return cleanedTournament, nil
}

func ConvertTimeToString(startTime time.Time) string {
	hour := startTime.Hour()
	minute := startTime.Minute()
	switch {
	case hour < 10 && minute < 10:
		return fmt.Sprintf("0%v0%v", hour, minute)
	case hour < 10:
		return fmt.Sprintf("0%v%v", hour, minute)
	case minute < 10:
		return fmt.Sprintf("%v0%v", hour, minute)
	default:
		return fmt.Sprintf("%v%v", hour, minute)
	}
}

func ConvertDateToString(startTime time.Time) string {
	day := startTime.Day()
	month := startTime.Month().String()[:3]
	if day < 10 {
		return fmt.Sprintf("0%v-%v", day, month)
	}
	return fmt.Sprintf("%v-%v", day, month)
}
