package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func formatHM(d time.Duration) string {
	totalMinutes := int(d / time.Minute)
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func LoadHistort(filename string, dayRef int) (Today string, ThisWeek string, LastWeek string, err error) {
	Today = "00:00"
	ThisWeek = "00:00"
	LastWeek = "00:00"

	file, err := os.Open(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return Today, ThisWeek, LastWeek, nil
		}
		return Today, ThisWeek, LastWeek, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return Today, ThisWeek, LastWeek, fmt.Errorf("failed to read CSV: %w", err)
	}

	// Time windows
	now := time.Now()
	loc := now.Location()

	// Start of today
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	// Start of current week (Monday)
	weekday := int(todayStart.Weekday())
	// Go's Weekday: Sunday=0 ... Saturday=6
	daysSinceRef := (weekday + 7 - dayRef) % 7

	thisWeekStart := todayStart.AddDate(0, 0, -daysSinceRef)
	thisWeekEnd := thisWeekStart.AddDate(0, 0, 7)
	lastWeekStart := thisWeekStart.AddDate(0, 0, -7)

	var todayDur, thisWeekDur, lastWeekDur time.Duration

	for i := 1; i < len(records); i++ {
		rec := records[i]
		if len(rec) < 2 {
			continue
		}
		startStr := rec[0]
		endStr := rec[1]

		// Skip incomplete or zero-value rows
		if startStr == "" || endStr == "" {
			fmt.Println("incomplete line found")
			continue
		}

		startTime, errStart := time.ParseInLocation(layout, startStr, loc)
		endTime, errEnd := time.ParseInLocation(layout, endStr, loc)
		if errStart != nil || errEnd != nil || endTime.Before(startTime) {
			fmt.Println("invalid field found")
			continue
		}

		dur := endTime.Sub(startTime)
		if !endTime.Before(lastWeekStart) && endTime.Before(thisWeekStart) {
			lastWeekDur += dur
			continue
		}
		if !endTime.Before(todayStart) {
			todayDur += dur
		}
		if !endTime.Before(thisWeekStart) && endTime.Before(thisWeekEnd) {
			thisWeekDur += dur
		}

	}

	Today = formatHM(todayDur)
	ThisWeek = formatHM(thisWeekDur)
	LastWeek = formatHM(lastWeekDur)
	return Today, ThisWeek, LastWeek, nil
}

func AddStartToCSV(filename string, start time.Time) error {
	_, err := os.Stat(filename)
	newFile := os.IsNotExist(err)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	if newFile {
		if _, err := writer.WriteString("Start,End\n"); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
	}

	if _, err := writer.WriteString(fmt.Sprintf("%s", start.Format(layout))); err != nil {
		return fmt.Errorf("failed to write CSV: %w", err)
	}

	return nil
}

func AddEndToCSV(filename string, end time.Time) error {
	_, err := os.Stat(filename)
	newFile := os.IsNotExist(err)

	if newFile {
		return fmt.Errorf("file should exist: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	if _, err := writer.WriteString(fmt.Sprintf(",%s\n", end.Format(layout))); err != nil {
		return fmt.Errorf("failed to write CSV: %w", err)
	}

	return nil
}

func GetLastTime(filename string) (string, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", fmt.Errorf("file does not exist")
		}
		return "", "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return "", "", fmt.Errorf("failed to read CSV: %w", err)
	}

	// If no data records (only header), return zero time
	if len(records) <= 1 {
		return "", "", fmt.Errorf("no time entries found")
	}

	// Get the last record
	lastRecord := records[len(records)-1]
	if len(lastRecord) < 2 {
		return "", "", fmt.Errorf("incomplete last record")
	}

	startStr := lastRecord[0]
	endStr := lastRecord[1]

	return startStr, endStr, nil
}

func ReplaceLastRecord(filename string, startStr, endStr string) error {
	// Parse the input strings to validate them
	loc := time.Now().Location()
	_, err := time.ParseInLocation(layout, startStr, loc)
	if err != nil {
		return fmt.Errorf("failed to parse start time: %w", err)
	}
	_, err = time.ParseInLocation(layout, endStr, loc)
	if err != nil {
		return fmt.Errorf("failed to parse end time: %w", err)
	}

	// Read all records
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	// Check if we have any data records
	if len(records) <= 1 {
		return fmt.Errorf("no time entries found to replace")
	}

	// Replace the last record
	records[len(records)-1] = []string{startStr, endStr}

	// Write back to file
	file, err = os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}
