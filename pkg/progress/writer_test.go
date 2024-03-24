package progress

import (
	"testing"
)

// MockWriter is a mock implementation of the io.Writer interface for testing.
type MockWriter struct {
	writeFunc func(p []byte) (n int, err error)
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	if m.writeFunc != nil {
		return m.writeFunc(p)
	}
	return 0, nil
}

// Test cases for the Writer methods.
var writerTests = []struct {
	name        string
	writeFunc   func(p []byte) (n int, err error)
	expectedN   int64
	expectedErr error
}{
	{"Write Success", func(p []byte) (n int, err error) { return len(p), nil }, 5, nil},
	// {"Write Error", func(p []byte) (n int, err error) { return 0, errors.New("test error") }, 0, errors.New("test error")},
}

func TestWriter(t *testing.T) {
	for _, tc := range writerTests {
		t.Run(tc.name, func(t *testing.T) {
			mockWriter := &MockWriter{writeFunc: tc.writeFunc}
			writer := NewWriter(mockWriter)

			data := []byte("Hello")
			n, err := writer.Write(data)

			if n != len(data) {
				t.Errorf("Expected Write() to write %d bytes, got %d", len(data), n)
			}

			if err != tc.expectedErr {
				t.Errorf("Expected error '%v', got '%v'", tc.expectedErr, err)
			}

			if writer.N() != tc.expectedN {
				t.Errorf("Expected N() to return %d, got %d", tc.expectedN, writer.N())
			}

			if writer.Err() != tc.expectedErr {
				t.Errorf("Expected Err() to return '%v', got '%v'", tc.expectedErr, writer.Err())
			}
		})
	}
}
