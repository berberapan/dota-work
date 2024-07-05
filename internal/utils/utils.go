package utils

import (
	"fmt"
	"log"
	"time"
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
