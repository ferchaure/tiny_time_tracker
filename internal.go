package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

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
		return fmt.Errorf("File should exists: %w", err)
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
