package tools

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

// An in-memory log store for testing log outputs.

type TestLogBuffer struct {
	GotBuffer    *bytes.Buffer
	ExpectBuffer *bytes.Buffer
}

// Add a log string to the list of expected log outputs.
func (b TestLogBuffer) ExpectLog(s string) error {
	_, err := b.ExpectBuffer.WriteString(s)
	if err != nil {
		return err
	}
	_, err = b.ExpectBuffer.WriteString("\n")
	if err != nil {
		return err
	}
	return nil
}

// Tests the received log against the expectations.
func (b *TestLogBuffer) TestLogValues(t *testing.T) {
	if strings.Compare(b.GotBuffer.String(), b.ExpectBuffer.String()) != 0 {
		gotLog := b.GotBuffer.String()
		wantLog := b.ExpectBuffer.String()
		t.Errorf("Got wrong log output. Got:\n%+s Want:\n%+s", gotLog, wantLog)
	}
}

// CreateAndActivateEmptyTestLogBuffer creates a new TestLogBuffer and applies it to the standard log library.
func CreateAndActivateEmptyTestLogBuffer() *TestLogBuffer {
	logBuffer := TestLogBuffer{
		&bytes.Buffer{},
		&bytes.Buffer{},
	}
	log.SetOutput(logBuffer.GotBuffer)
	log.SetFlags(0)
	return &logBuffer
}
