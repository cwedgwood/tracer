//

package tracer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

func newTestLogger(b *bytes.Buffer) logr.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zl := zerolog.New(b).With().Caller().Timestamp().Logger()
	return zerologr.New(&zl)
}

type logOutput struct {
	Level       string    `json:"level"`                 // error or a 'v level';  0:"info", 1:"debug", 2:"trace"
	TraceId     string    `json:"traceid"`               //
	TraceOrigin string    `json:"traceorigin,omitempty"` //
	V           *int      `json:"v,omitempty"`           // present for non-errors, 0 through 2 see Level
	Error       string    `json:"error,omitempty"`       // present for errors, error string
	Caller      string    `json:"caller"`                //
	Time        time.Time `json:"time"`                  //
	Message     string    `json:"message"`               //
}

func TestLoggerInfoFields(t *testing.T) {
	logBuffer := &bytes.Buffer{}
	ctx := ContextLoggerWithTraceId(context.TODO(), newTestLogger(logBuffer), GenerateTraceId, "info-tester")

	// this is how most users of the context will consume it and use it
	log := logr.FromContextOrDiscard(ctx)
	log.Info("something1")

	var logOut logOutput
	err := json.Unmarshal(logBuffer.Bytes(), &logOut)
	if err != nil {
		t.Errorf("Failed to unmarshal log output: %s", err)
	}

	if logOut.Level != "info" {
		t.Errorf("Expected info, got %s", logOut.Level)
	}
	if logOut.TraceId == "" {
		t.Errorf("Expected non-empty traceid")
	}
	if logOut.TraceOrigin != "info-tester" {
		t.Errorf("Expected tester, got %s", logOut.TraceOrigin)
	}
	if logOut.V != nil && *logOut.V != 0 {
		t.Errorf("Expected 0, got %d", *logOut.V)
	}
	if logOut.Caller == "" {
		t.Errorf("Expected non-empty caller")
	}
	if logOut.Time.IsZero() || time.Since(logOut.Time) > time.Second {
		t.Errorf("Expected non-zero time")
	}
	if logOut.Message != "something1" {
		t.Errorf("Expected something1, got %s", logOut.Message)
	}
}

func TestLoggerError(t *testing.T) {
	logBuffer := &bytes.Buffer{}
	ctx := ContextLoggerWithTraceId(context.TODO(), newTestLogger(logBuffer), GenerateTraceId, "error-tester")

	// this is how most users of the context will consume it and use it
	log := logr.FromContextOrDiscard(ctx)
	log.Error(errors.New("testerror"), "error testing")

	var logOut logOutput
	err := json.Unmarshal(logBuffer.Bytes(), &logOut)
	if err != nil {
		t.Errorf("Failed to unmarshal log output: %s", err)
	}
	if logOut.Level != "error" {
		t.Errorf("Expected error, got %s", logOut.Level)
	}
	if logOut.Error != "testerror" {
		t.Errorf("Expected testerror, got %s", logOut.Error)
	}
	if logOut.Message != "error testing" {
		t.Errorf("Expected error testing, got %s", logOut.Message)
	}
}

// TestPresentTraceIdAndOrigin tests that the traceid and traceorigin are present in the context.
func TestPresentTraceIdAndOrigin(t *testing.T) {
	ctx := ContextLoggerWithTraceId(context.TODO(), newTestLogger(&bytes.Buffer{}), GenerateTraceId, "origin.tester")
	id, origin := TraceIdAndOrigin(ctx)
	if id == "" {
		t.Errorf("Expected non-empty traceid")
	}
	if origin != "origin.tester" {
		t.Errorf("Expected origin.tester, got %s", origin)
	}
}

// TestMissingTraceIdAndOrigin tests that the traceid and traceorigin are not present in the context and that we do not panic.
func TestMissingTraceIdAndOrigin(t *testing.T) {
	ctx := context.TODO()
	id, origin := TraceIdAndOrigin(ctx)
	if len(id) != 0 {
		t.Errorf("Expected empty traceid, got %s", id)
	}
	if len(origin) != 0 {
		t.Errorf("Expected empty traceorigin, got %s", origin)
	}
}
