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

// Package downloader implements a high-level interface for downloading
// a Factorio server and installing it.
package downloader

import (
	"context"
	"fmt"
	"os"

	"github.com/jaredallard/factorio-docker/internal/factorio"
)

// Download downloads the specified Factorio version to the output
// directory. If version is "stable" or "experimental", it will download
// the latest stable or experimental version, respectively.
func Download(ctx context.Context, version, sha256sum, outputDir string) error {
	// Ensure the output directory is empty.
	if files, err := os.ReadDir(outputDir); err == nil && len(files) > 0 {
		return fmt.Errorf("output directory %s is not empty", outputDir)
	}

	// Special case, resolve stable/experimental to the latest version of
	// their respective channels.
	if VersionIsChannel(version) {
		var err error
		version, err = GetVersionForChannel(Channel(version))
		if err != nil {
			return err
		}
	}

	// If there's no SHA256 sum, get it.
	if sha256sum == "" {
		var err error
		sha256sum, err = factorio.GetSHA256(version)
		if err != nil {
			return err
		}
	}

	return factorio.DownloadVersion(ctx, version, sha256sum, outputDir)
}
