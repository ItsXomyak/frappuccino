package helper

import (
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"
)

func IsValidName(name string) error {
	nameLength := len(name)
	if nameLength < 2 || nameLength > 120 {
		return fmt.Errorf("the name must be between 2 and 63 characters long: %s", name)
	}

	validNameRegex := regexp.MustCompile(`^[A-Za-z0-9_ ]+$`)

	if !validNameRegex.MatchString(name) {
		return fmt.Errorf("name must only contain alphanumeric characters, underscores, and hyphens: %s", name)
	}

	if strings.Contains(name, "..") {
		return fmt.Errorf("Name cannot contain consecutive periods")
	}
	if strings.Contains(name, "--") {
		return fmt.Errorf("Name cannot contain consecutive dashes")
	}
	if strings.Contains(name, "/") {
		return fmt.Errorf("Name cannot contain consecutive slash")
	}
	return nil
}

func HelperFunc() {
	fmt.Printf(`
	Usage:
  frappuccino [--port <N>] [--dir <S>] 
  frappuccino --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --dir S      Path to the data directory.
	`)
}

var counter int64 = 0

func GenerateID() (string, error) {
	newValue := atomic.AddInt64(&counter, 1)
	if newValue < 1 {
		return "", fmt.Errorf("counter has overflowed")
	}
	return fmt.Sprintf("%d", newValue), nil
}

func ClearId(id string) string {
	cleanStr := strings.ReplaceAll(id, `\`, "")
	cleanStr = strings.ReplaceAll(cleanStr, `"`, "")
	return cleanStr
}
