package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetFileName(date time.Time) string {
	return fmt.Sprintf("data_%04d_%02d_%02d.json", date.Year(), date.Month(), date.Day())
}

func GetFileDate(fileName string) (*time.Time, error) {
	fileParts := strings.Split(fileName, ".")
	parts := strings.Split(fileParts[0], "_")
	year, _ := strconv.Atoi(parts[1])
	month, _ := strconv.Atoi(parts[2])
	day, _ := strconv.Atoi(parts[3])
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return &date, nil
}

func GetFolder(date time.Time) string {
	return fmt.Sprintf("%04d-%02d/", date.Year(), date.Month())
}

func GetWholePath(date time.Time) string {
	return GetFolder(date) + GetFileName(date)
}
