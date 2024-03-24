package progress

import (
	"bytes"
	"testing"
)

func TestReader_Read(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		readSize    int
		expectedN   int
		expectedErr error
	}{
		{"Read Success", []byte("Hello, world!"), 5, 5, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a buffer with test data
			buffer := bytes.NewBuffer(tc.data)
			// Create a new Reader with the buffer
			reader := NewReader(buffer)
			// Read from the reader
			readData := make([]byte, tc.readSize)
			n, err := reader.Read(readData)
			// Verify the number of bytes read
			if n != tc.expectedN {
				t.Errorf("Expected to read %d bytes, got %d", tc.expectedN, n)
			}
			// Verify the error
			if err != tc.expectedErr {
				t.Errorf("Expected error '%v', got '%v'", tc.expectedErr, err)
			}
		})
	}
}

func TestReader_N(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		readSize  int
		expectedN int64
	}{
		{"N Success", []byte("Hello, world!"), 5, 5},
		{"N After Error", []byte("Hello, world!"), 100, 13},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a buffer with test data
			buffer := bytes.NewBuffer(tc.data)
			// Create a new Reader with the buffer
			reader := NewReader(buffer)
			// Read from the reader
			_ = make([]byte, tc.readSize)
			// Verify the number of bytes read
			if reader.N() != tc.expectedN {
				t.Errorf("Expected N() to return %d, got %d", tc.expectedN, reader.N())
			}
		})
	}
}

func TestReader_Err(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		readSize       int
		expectedErrNil bool
	}{
		{"Err Success", []byte("Hello, world!"), 5, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a buffer with test data
			buffer := bytes.NewBuffer(tc.data)
			// Create a new Reader with the buffer
			reader := NewReader(buffer)
			// Read from the reader
			_ = make([]byte, tc.readSize)
			// Verify the error
			if tc.expectedErrNil && reader.Err() != nil {
				t.Errorf("Expected Err() to return nil, got %v", reader.Err())
			} else if !tc.expectedErrNil && reader.Err() == nil {
				t.Errorf("Expected Err() to return an error, got nil")
			}
		})
	}
}
