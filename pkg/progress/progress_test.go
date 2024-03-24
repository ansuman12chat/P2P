package progress

import (
	"context"
	"errors"
	"testing"
	"time"
)

// MockCounter is a mock implementation of the Counter interface for testing.
type MockCounter struct {
	count int64
	err   error
}

func (m *MockCounter) N() int64 {
	return m.count
}

func (m *MockCounter) Err() error {
	return m.err
}

func TestProgress_N(t *testing.T) {
	// Create a Progress struct with a specific count
	p := Progress{n: 100}
	// Test the N method
	if p.N() != 100 {
		t.Errorf("Expected N() to return %d, got %d", 100, p.N())
	}
}

func TestProgress_Size(t *testing.T) {
	// Create a Progress struct with a specific size
	p := Progress{size: 200}
	// Test the Size method
	if p.Size() != 200 {
		t.Errorf("Expected Size() to return %d, got %d", 200, p.Size())
	}
}

func TestProgress_Complete(t *testing.T) {
	// Create a Progress struct with an error set to io.EOF
	p := Progress{err: errors.New("EOF")}
	// Test the Complete method
	if !p.Complete() {
		t.Errorf("Expected Complete() to return true for io.EOF error, got false")
	}
}

func TestProgress_Percent(t *testing.T) {
	// Create a Progress struct with specific count and size
	p := Progress{n: 50, size: 200}
	// Test the Percent method
	if p.Percent() != 25 {
		t.Errorf("Expected Percent() to return %f, got %f", 25.0, p.Percent())
	}
}

func TestProgress_Remaining(t *testing.T) {
	// Create a Progress struct with a specific estimated time
	p := Progress{estimated: time.Now().Add(time.Hour)}
	// Test the Remaining method
	if p.Remaining() == -1 {
		t.Errorf("Expected Remaining() to return a non-negative duration, got -1")
	}
}

func TestProgress_Estimated(t *testing.T) {
	// Create a Progress struct with a specific estimated time
	estTime := time.Now().Add(time.Hour)
	p := Progress{estimated: estTime}
	// Test the Estimated method
	if p.Estimated() != estTime {
		t.Errorf("Expected Estimated() to return %v, got %v", estTime, p.Estimated())
	}
}

func TestNewTicker(t *testing.T) {
	// Create a mock counter
	counter := &MockCounter{}
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Call NewTicker with the mock counter and context
	ch := NewTicker(ctx, counter, 100, time.Millisecond)
	// Loop over the channel to receive progress updates
	for progress := range ch {
		// Verify the progress updates as needed
		t.Logf("Progress: %v", progress)
	}
}
