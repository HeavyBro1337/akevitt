package akevitt

import (
	"fmt"
	"time"
)

var format string = "[%s] %s: %s"

func LogInfo(message string) {
	fmt.Printf(format, time.Now(), "LOG", message)
}
func LogWarn(message string) {
	fmt.Printf(format, time.Now(), "WARN", message)
}
func LogError(message string) {
	fmt.Printf(format, time.Now(), "ERR", message)
}
