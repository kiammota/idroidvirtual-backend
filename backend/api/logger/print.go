package logger

import (
	"fmt"
	"strings"
	"time"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

func log(color string, level string, message ...string) {
	timestamp := time.Now().Format("2006/01/02 - 15:04:05")
	fullMessage := strings.Join(message, " ")

	fmt.Printf(
		"%s[%s] | %-5s | %s%s\n",
		color,
		timestamp,
		level,
		fullMessage,
		Reset,
	)
}

func OK(message ...string) {
	log(Green, "OK", message...)
}

func Err(message ...string) {
	log(Red, "ERROR", message...)
}

func Warn(message ...string) {
	log(Yellow, "WARN", message...)
}

func Log(message ...string) {
	log(Cyan, "INFO", message...)
}
