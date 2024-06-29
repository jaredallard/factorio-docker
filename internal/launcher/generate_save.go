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

package launcher

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jaredallard/factorio-docker/internal/config"
)

// generateSave generates a new Factorio save.
func generateSave(cfg *config.Config, execPath, saveName string) error {
	args := [...]string{
		execPath,
		"--create", filepath.Join(cfg.ServerDataPath, "saves", saveName) + ".zip",
		"--map-gen-settings", filepath.Join(cfg.ServerDataPath, "map-gen-settings.json"),
		"--map-settings", filepath.Join(cfg.ServerDataPath, "map-settings.json"),
	}

	//nolint:gosec // Why: created above
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = cfg.ServerDataPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateDefaultSave(cfg *config.Config, execPath string) error {
	// If there's no save found, create one.
	savesDir := filepath.Join(cfg.ServerDataPath, "saves")
	files, err := os.ReadDir(savesDir)
	if err != nil {
		return err
	}

	// Check if there's any zip files.
	var foundSave bool
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".zip" {
			foundSave = true
			break
		}
	}
	if foundSave {
		return nil
	}

	return generateSave(cfg, execPath, "_autosave1")
}
