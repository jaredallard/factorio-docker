package factorio_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaredallard/factorio-docker/internal/factorio"
	"gotest.tools/v3/assert"
)

func TestCanGetLatestAndSHA256(t *testing.T) {
	rels, err := factorio.GetLatestReleases()
	assert.NilError(t, err)

	// Ensure we have a stable and experimental release
	assert.Assert(t, rels.Stable != "")
	assert.Assert(t, rels.Experimental != "")

	// Ensure we can get the SHA256 of the stable release
	stableSHA, err := factorio.GetSHA256(rels.Stable)
	assert.NilError(t, err)
	assert.Assert(t, stableSHA != "")

	// Ensure we can get the SHA256 of the experimental release
	expSHA, err := factorio.GetSHA256(rels.Experimental)
	assert.NilError(t, err)
	assert.Assert(t, expSHA != "")

	// Ensure we can download the stable version
	tmpDir := t.TempDir()

	t.Log("Downloading", rels.Stable, "to", tmpDir, "sha256:", stableSHA)
	err = factorio.DownloadVersion(rels.Stable, stableSHA, tmpDir)
	assert.NilError(t, err)

	// Ensure it contains an expected file
	_, err = os.Stat(filepath.Join(tmpDir, "factorio", "bin", "x64", "factorio"))
	assert.NilError(t, err)
}
