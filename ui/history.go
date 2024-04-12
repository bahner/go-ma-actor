package ui

import (
	"bufio"
	"fmt"
	"os"

	"strings"

	"github.com/bahner/go-ma-actor/config"
	"github.com/spf13/viper"
)

func historySize() int {
	return viper.GetInt("ui.history-size")
}

func (ui *ChatUI) pushToHistory(line string) {

	historySize := historySize()

	if len(ui.inputHistory) == historySize {
		// Remove the oldest entry when we reach max size
		copy(ui.inputHistory, ui.inputHistory[1:])
		ui.inputHistory = ui.inputHistory[:historySize-1]
	}
	ui.inputHistory = append(ui.inputHistory, line)
	if err := appendToPersistentHistory(line + "\n"); err != nil {
		fmt.Printf("Error appending to history file: %v\n", err)
	}
}

// appendToFile opens the specified file in append mode and writes the string to it.
func appendToPersistentHistory(text string) error {

	if !(strings.HasPrefix(text, "/") ||
		strings.HasPrefix(text, ".")) {
		return nil
	}

	filename := config.DBHistory()
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(text); err != nil {
		return err
	}

	return nil
}

// loadHistory loads the last 'historySize' lines from the history file into the input history.
func (ui *ChatUI) loadHistory() error {
	filename := config.DBHistory()

	// Open the file for reading.
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// It's okay if the file doesn't exist yet.
			return nil
		}
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	// Read all lines into 'lines'. This is not memory efficient for large files but is simple.
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if len(lines) > historySize() {
		lines = lines[len(lines)-historySize():]
	}

	ui.inputHistory = lines
	return nil
}
