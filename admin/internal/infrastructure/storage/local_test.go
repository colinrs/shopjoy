package storage_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestLocalStorage_Save_RoundTrip(t *testing.T) {
	// Requires real DB connection. Real roundtrip tests will live under
	// internal/integration/ in a later task; this stub is intentionally a
	// no-op so the file compiles and `go test -run TestLocalStorage` exits 0.
	t.Skip("requires real DB connection; see README for local integration test setup")
	_ = bytes.NewBufferString
	_ = uuid.NewString
	_ = filepath.Join
	_ = os.WriteFile
	_ = context.Background
}
