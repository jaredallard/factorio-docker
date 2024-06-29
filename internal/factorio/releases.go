// Copyright (C) 2024 Jared Allard
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package factorio

import (
	"archive/tar"
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/ulikunitz/xz"
)

// DownloadVersion downloads a Factorio version to the specified
// directory. The downloaded version is validated against the SHA256 sum
// on the remote, and extracted to the specified directory.
func DownloadVersion(ctx context.Context, version, sha256sum, destDir string) error {
	if _, err := os.Stat(destDir); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx,
		"GET", fmt.Sprintf("https://factorio.com/get-download/%s/headless/linux64", version),
		http.NoBody)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	h := sha256.New()
	// While extracting the file, also calculate the SHA256sum. This
	// prevents needing to read the file twice.
	xzr, err := xz.NewReader(bufio.NewReader(io.TeeReader(resp.Body, h)))
	if err != nil {
		return err
	}

	// Extract the tarball to the destination directory.
	t := tar.NewReader(xzr)
	for {
		header, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fInf := header.FileInfo()

		archiveFilePath := header.Name
		// Remove factorio/ from the path.
		archiveFilePath = filepath.Clean(strings.TrimPrefix(archiveFilePath, "factorio/"))

		//nolint:gosec // Why: We're mitigating it below.
		destpath := filepath.Join(destDir, archiveFilePath)
		if !strings.HasPrefix(destpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", archiveFilePath)
		}

		if header.Typeflag == tar.TypeDir {
			// We JIT create them later.
			continue
		}

		// Using a closure to ensure the file is closed after we're done
		// since defer won't work in a loop.
		if err := func() error {
			if err := os.MkdirAll(filepath.Dir(destpath), 0o755); err != nil {
				return err
			}

			f, err := os.Create(destpath)
			if err != nil {
				return err
			}
			defer f.Close()

			if err := os.Chmod(destpath, fInf.Mode()); err != nil {
				return fmt.Errorf("failed to chmod file %s: %w", destpath, err)
			}

			fmt.Println(" ->", destpath, humanize.Bytes(uint64(header.Size)))
			for {
				_, err = io.CopyN(f, t, 4096)
				if err != nil {
					if errors.Is(err, io.EOF) {
						err = nil
					}

					break
				}
			}
			// Pass the error up the stack.
			return err
		}(); err != nil {
			return err
		}
	}

	// Validate the SHA256 sum.
	hexSum := hex.EncodeToString(h.Sum(nil))
	if hexSum != sha256sum {
		return fmt.Errorf("SHA256 sum does not match: expected %s, got %s", sha256sum, hexSum)
	}

	fmt.Println("Downloaded and validated Factorio version", version)

	return nil
}
