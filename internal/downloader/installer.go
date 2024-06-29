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

package downloader

import (
	"context"
	_ "embed" // Embed the default Factorio config.
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jaredallard/factorio-docker/internal/config"
	"github.com/jaredallard/factorio-docker/internal/state"
)

// EnsureVersion ensures that the Factorio server is installed and up-to-date.
func EnsureVersion(ctx context.Context, cfg *config.Config, log *slog.Logger) error {
	st := state.Open(filepath.Join(cfg.InstallPath, "state.json"))

	if st.Version == "" {
		log.Info("No version of Factorio installed, installing...")

		// Ensure the install path exists.
		if _, err := os.Stat(cfg.InstallPath); os.IsNotExist(err) {
			if err := os.MkdirAll(cfg.InstallPath, 0o755); err != nil {
				return fmt.Errorf("failed to create install directory: %w", err)
			}
		}
	}

	// If we're installing a channel, resolve it to the latest version.
	if VersionIsChannel(cfg.Version) {
		// If the version is a channel, resolve it to the latest version.
		ver, err := GetVersionForChannel(Channel(cfg.Version))
		if err != nil {
			return err
		}

		oldVer := cfg.Version
		cfg.Version = ver

		log.Info("Resolved channel", "channel", oldVer, "version", ver)
	}

	// If the installed version is the requested version, we're done.
	if st.Version == cfg.Version {
		log.Info("Factorio is desired version")
		return nil
	}

	log.Info("Installing Factorio", "version", cfg.Version, "previous", st.Version)

	// Ensure the directory is empty.
	if files, err := os.ReadDir(cfg.InstallPath); err == nil && len(files) > 0 {
		for _, f := range files {
			if err := os.RemoveAll(filepath.Join(cfg.InstallPath, f.Name())); err != nil {
				return fmt.Errorf("failed to remove existing Factorio installation: %w", err)
			}
		}
	}

	// The installed version is not the requested version, download it.
	if err := Download(ctx, cfg.Version, "", cfg.InstallPath); err != nil {
		return err
	}

	st.Version = cfg.Version
	if err := st.Save(); err != nil {
		return fmt.Errorf("failed to track installed version: %w", err)
	}

	return nil
}
