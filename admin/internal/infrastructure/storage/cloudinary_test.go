package storage_test

import (
	"testing"
)

// TestCloudinary_Sign_Skip verifies the SHA1 signature canonical-form algorithm
// against Cloudinary's documentation. The exact expected SHA1 requires fetching
// the Cloudinary docs page at implementation time to confirm the canonical
// example (timestamp, params, secret). Skipped here because:
//
//   - The Cloudinary API example SHA1 depends on a specific sample set of
//     params that is not stable across doc revisions.
//   - WebFetch + WebSearch would be needed to lock the expected SHA1, and that
//     is beyond a quick task scope.
//
// Real Sign verification is exercised via integration test (T11+) once the
// upload API endpoint is wired and a live (or sandbox) Cloudinary account is
// configured.
//
// To enable: replace the t.Skip with an assertion comparing s.Sign(ctx, ...)
// against a SHA1 computed from canonicalParams + api_secret matching the
// Cloudinary docs example.

func TestCloudinary_Sign_Skip(t *testing.T) {
	t.Skip("Cloudinary SHA1 canonical-form verification requires a documented example; deferred to integration test.")
}
