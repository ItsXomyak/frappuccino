package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func LoadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
		return
	}
	defer file.Close()

	var lines []string
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			break
		}
		lines = append(lines, strings.Split(string(buf[:n]), "\n")...)
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}
}

func GetEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
