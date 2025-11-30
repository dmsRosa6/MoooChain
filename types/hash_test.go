package types

import (
	"testing"
)

func TestHashFromBytes_Valid(t *testing.T) {
	input := make([]byte, 32)
	for i := range input {
		input[i] = byte(i)
	}

	h := HashFromBytes(input)

	for i := range input {
		if h[i] != input[i] {
			t.Fatalf("expected h[%d] = %d, got %d", i, input[i], h[i])
		}
	}
}

func TestHashFromBytes_InvalidSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for invalid size, got none")
		}
	}()

	_ = HashFromBytes([]byte{1, 2, 3})
}

func TestRandomBytes_Size(t *testing.T) {
	b := RandomBytes(16)
	if len(b) != 16 {
		t.Fatalf("expected 16 bytes, got %d", len(b))
	}
}

func TestRandomBytes_InvalidSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for negative size")
		}
	}()

	_ = RandomBytes(-1)
}

func TestRandomHash(t *testing.T) {
	h := RandomHash()

	if len(h) != 32 {
		t.Fatalf("expected hash of size 32, got %d", len(h))
	}

	// Should not be all zeros
	zero := true
	for _, b := range h {
		if b != 0 {
			zero = false
			break
		}
	}

	if zero {
		t.Fatalf("expected random hash, not zeroed hash")
	}
}
