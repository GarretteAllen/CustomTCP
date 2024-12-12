package utils

import "log"

func LogError(msg string, err error) {
	log.Printf("ERROR: %s: %v", msg, err)
}

func LogInfo(format string, args ...interface{}) {
	if len(args) > 0 {
		format = format + " %v"
	}
	log.Printf("INFO: "+format, args...)
}
