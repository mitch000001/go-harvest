package harvest

import (
	"io"
	"log"
	"os"
)

var debugMode bool = false
var infoMode bool = true
var debugFlags int = log.Lshortfile | log.LstdFlags

func DebugMode() bool {
	return debugMode
}

func SetDebugMode(debug bool) {
	debugMode = debug
}

func InfoLog() bool {
	return infoMode
}

func SetInfoLog(info bool) {
	infoMode = info
}

func Debug(fn func()) {
	debugMode = true
	fn()
	debugMode = false
}

var debug *log.Logger = newConditionalLogger(os.Stdout, "harvest: ", debugFlags, &debugMode)
var info *log.Logger = newConditionalLogger(os.Stdout, "harvest: ", log.LstdFlags, &infoMode)

func newConditionalLogger(w io.Writer, prefix string, flag int, condition *bool) *log.Logger {
	condWriter := newConditionalWriter(w, condition)
	return log.New(condWriter, prefix, flag)
}

func newConditionalWriter(w io.Writer, condition *bool) io.Writer {
	return &conditionalWriter{Writer: w, condition: condition}
}

type conditionalWriter struct {
	condition *bool
	io.Writer
}

func (c *conditionalWriter) Write(p []byte) (n int, err error) {
	if *c.condition {
		return c.Writer.Write(p)
	}
	return 0, nil
}
