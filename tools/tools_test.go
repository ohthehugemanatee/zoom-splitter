package tools

import (
	"log"
	"testing"
)

func TestTestLogBuffer(t *testing.T) {
	line1 := "One is the loneliest number"
	line2 := "It takes two to tango"
	line3 := "Three's a crowd"
	t.Run("Identical single lines should pass", func(t *testing.T) {
		logBuffer := CreateAndActivateEmptyTestLogBuffer()
		logBuffer.ExpectLog(line1)
		log.Print(line1)
		newT := testing.T{}
		logBuffer.TestLogValues(&newT)
		if newT.Failed() {
			t.Error("Identical log lines recorded a test failure")
		}
	})
	t.Run("Different single lines should fail", func(t *testing.T) {
		logBuffer := CreateAndActivateEmptyTestLogBuffer()
		logBuffer.ExpectLog(line1)
		log.Print(line2)
		newT := testing.T{}
		logBuffer.TestLogValues(&newT)
		if !newT.Failed() {
			t.Error("Different log lines did not result in test failure")
		}
	})
	t.Run("Identical multiple lines should pass", func(t *testing.T) {
		logBuffer := CreateAndActivateEmptyTestLogBuffer()
		logBuffer.ExpectLog(line1)
		logBuffer.ExpectLog(line2)
		logBuffer.ExpectLog(line3)
		log.Print(line1)
		log.Print(line2)
		log.Print(line3)
		newT := testing.T{}
		logBuffer.TestLogValues(&newT)
		if newT.Failed() {
			t.Error("Identical multiple log lines recorded as a test failure")
		}
	})
	t.Run("Missing log lines should fail", func(t *testing.T) {
		logBuffer := CreateAndActivateEmptyTestLogBuffer()
		logBuffer.ExpectLog(line1)
		logBuffer.ExpectLog(line2)
		logBuffer.ExpectLog(line3)
		log.Print(line3)
		log.Print(line1)
		newT := testing.T{}
		logBuffer.TestLogValues(&newT)
		if !newT.Failed() {
			t.Error("Missing log lines did not result in test failure")
		}
	})
	t.Run("Extra unexepcted log lines should fail", func(t *testing.T) {
		logBuffer := CreateAndActivateEmptyTestLogBuffer()
		logBuffer.ExpectLog(line1)
		logBuffer.ExpectLog(line3)
		log.Print(line3)
		log.Print(line2)
		log.Print(line1)
		newT := testing.T{}
		logBuffer.TestLogValues(&newT)
		if !newT.Failed() {
			t.Error("Extra unexpected log lines did not result in test failure")
		}
	})
	t.Run("Identical multiple lines in different order should fail", func(t *testing.T) {
		logBuffer := CreateAndActivateEmptyTestLogBuffer()
		logBuffer.ExpectLog(line1)
		logBuffer.ExpectLog(line2)
		logBuffer.ExpectLog(line3)
		log.Print(line3)
		log.Print(line2)
		log.Print(line1)
		newT := testing.T{}
		logBuffer.TestLogValues(&newT)
		if !newT.Failed() {
			t.Error("Identical multiple log lines in the wrong order did not result in test failure")
		}
	})

}
