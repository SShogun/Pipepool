package testutil

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func NewBufferLogger(buf *bytes.Buffer) *slog.Logger {
	if buf == nil {
		return NewDiscardLogger()
	}

	return slog.New(slog.NewTextHandler(buf, nil))
}

func BuildManyInputs(n int, value string) []string {
	if n <= 0 {
		return nil
	}

	inputs := make([]string, n)
	for i := 0; i < n; i++ {
		inputs[i] = value
	}

	return inputs
}

func BuildInputs(values ...string) []string {
	inputs := make([]string, 0, len(values))
	for _, value := range values {
		inputs = append(inputs, strings.Clone(value))
	}

	return inputs
}
