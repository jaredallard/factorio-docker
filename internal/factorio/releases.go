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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jamespfennell/xz"
)

// DownloadVersion downloads a Factorio version to the specified
// directory. The downloaded version is validated against the SHA256 sum
// on the remote, and extracted to the specified directory.
func DownloadVersion(version, sha256sum, dest string) error {
	if _, err := os.Stat(dest); err != nil {
		return err
	}

	resp, err := http.Get(fmt.Sprintf("https://factorio.com/get-download/%s/headless/linux64", version))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	h := sha256.New()
	// While extracting the file, also calculate the SHA256sum. This
	// prevents needing to read the file twice.
	xzr := xz.NewReader(io.TeeReader(resp.Body, h))

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

		path := filepath.Join(dest, header.Name)
		if header.Typeflag == tar.TypeDir {
			// We JIT create them later.
			continue
		}

		// Using a closure to ensure the file is closed after we're done
		// since defer won't work in a loop.
		if err := func() error {
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}

			f, err := os.Create(path)
			if err != nil {
				return err
			}
			defer f.Close()

			if err := os.Chmod(path, fInf.Mode()); err != nil {
				return fmt.Errorf("failed to chmod file %s: %w", path, err)
			}

			fmt.Println("Extracting", header.Name, "->", path)
			_, err = io.Copy(f, t)
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
