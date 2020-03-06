package utils

import (
	"os"
	"strings"
)

// GoirateLogger acts as a logger for filtering certain types of problematic output.
// Main case is forms of "Unsolicited response received on idle HTTP channel", which Go
// for some reason thinks is a good idea to dump straight into stdout...
type GoirateLogger struct {
}

func (logger *GoirateLogger) Write(data []byte) (n int, err error) {

	str := string(data)

	if strings.Contains(str, "Unsolicited response received on idle HTTP channel") {
		return 0, nil
	}

	return os.Stdout.Write(data)
}
