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

func LoadHistort(filename string) (Today string, ThisWeek string, LastWeek string, err error) {
	Today = "00:00"
	ThisWeek = "00:00"
	LastWeek = "00:00"

	file, err := os.Open(filename)
	if err != nil {
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
	// Go's Weekday: Sunday=0 ... Saturday=6. We want Monday=0.
	daysSinceMonday := (weekday + 6) % 7
	thisWeekStart := todayStart.AddDate(0, 0, -daysSinceMonday)
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
