package debug

import "log"

var (
	Debug = false
)

func DEBUG(format string, args ...any) {
	if Debug {
		log.Printf(format, args...)
	}
}
