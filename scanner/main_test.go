package main

import (
	"bytes"
	"log"
	"os"
)

func CaptureCommand(cmd func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	cmd()
	log.SetOutput(os.Stdout)
	return buf.String()
}
