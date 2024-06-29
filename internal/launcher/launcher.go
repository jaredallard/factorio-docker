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

// Package launcher runs Factorio or wrappers, such as Factocord.
package launcher

import (
	_ "embed" // Embed is used to embed default configuration files.
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jaredallard/factorio-docker/internal/config"
	"github.com/jaredallard/factorio-docker/internal/factorio"
	"gopkg.in/ini.v1"
)

//go:embed embed/factorio-config.ini
var defaultFactorioConfig []byte

// setupConfig ensures that the Factorio server has a configuration file
// and that it points to our data directory.
func setupConfig(cfg *config.Config) error {
	configPath := filepath.Join(cfg.InstallPath, "config", "config.ini")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		if err := os.WriteFile(configPath, defaultFactorioConfig, 0o600); err != nil {
			return fmt.Errorf("failed to write default config: %w", err)
		}
	}

	// Disable pretty printing of the configuration file.
	ini.PrettyEqual = false
	ini.PrettyFormat = false

	conf, err := ini.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	// Point the server to the correct data path.
	if !conf.HasSection("path") {
		if _, err := conf.NewSection("path"); err != nil {
			return fmt.Errorf("failed to create path section: %w", err)
		}
	}
	if !conf.Section("path").HasKey("write-data") {
		if _, err := conf.Section("path").NewKey("write-data", cfg.ServerDataPath); err != nil {
			return fmt.Errorf("failed to create write-data key: %w", err)
		}
	} else {
		conf.Section("path").Key("write-data").SetValue(cfg.ServerDataPath)
	}

	if err := conf.SaveTo(configPath); err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}

	return nil
}

// installDefaultFiles installs default configuration files from the
// Factorio server install, if not present in the data directory.
func installDefaultFiles(log *slog.Logger, cfg *config.Config) error {
	type defaultFile struct {
		src, dest string
	}

	defaultDirs := []string{"mods", "saves", "scenarios"}
	for _, d := range defaultDirs {
		dirPath := filepath.Join(cfg.ServerDataPath, d)
		if _, err := os.Stat(dirPath); err != nil {
			log.Info("Creating default directory", "path", dirPath)
			if err := os.MkdirAll(dirPath, 0o755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		}
	}

	defaultFiles := []defaultFile{
		{
			src:  "server-settings.example.json",
			dest: "server-settings.json",
		},
		{
			src:  "map-gen-settings.example.json",
			dest: "map-gen-settings.json",
		},
		{
			src:  "map-settings.example.json",
			dest: "map-settings.json",
		},
	}

	for _, f := range defaultFiles {
		src := filepath.Join(cfg.InstallPath, "data", f.src)
		dest := filepath.Join(cfg.ServerDataPath, f.dest)

		if _, err := os.Stat(dest); err == nil {
			continue
		}

		log.Info("Installing default file", "src", src, "dest", dest)
		if err := func() error { // For defers in for loop.
			srcF, err := os.Open(src)
			if err != nil {
				return fmt.Errorf("failed to open source file: %w", err)
			}
			defer srcF.Close()

			destF, err := os.Create(dest)
			if err != nil {
				return fmt.Errorf("failed to create destination file: %w", err)
			}
			defer destF.Close()

			if _, err := io.Copy(destF, srcF); err != nil {
				return fmt.Errorf("failed to copy file: %w", err)
			}

			return nil
		}(); err != nil {
			return err
		}
	}

	return nil
}

// Launch starts a Factorio server based on the provided configuration.
func Launch(log *slog.Logger, cfg *config.Config) error {
	if err := setupConfig(cfg); err != nil {
		return fmt.Errorf("failed to setup config: %w", err)
	}

	if err := installDefaultFiles(log, cfg); err != nil {
		return fmt.Errorf("failed to install default files: %w", err)
	}

	execPath := filepath.Join(cfg.InstallPath, "bin", "x64", "factorio")

	if err := factorio.GenerateDefaultSave(cfg, execPath); err != nil {
		return fmt.Errorf("failed to generate default save: %w", err)
	}

	args := []string{
		execPath,

		// These settings allow us to keep the server files away from the
		// actual data.
		"--server-settings", filepath.Join(cfg.ServerDataPath, "server-settings.json"),
		"--server-banlist", filepath.Join(cfg.ServerDataPath, "server-banlist.json"),
		"--server-whitelist", filepath.Join(cfg.ServerDataPath, "server-whitelist.json"),
		"--server-adminlist", filepath.Join(cfg.ServerDataPath, "server-adminlist.json"),
		"--server-id", filepath.Join(cfg.ServerDataPath, "server-id.json"),
		"--use-server-whitelist",

		"--start-server-load-latest",
	}

	if cfg.Factocord.Enabled {
		return runFactocord(cfg, args)
	}

	return runVanilla(cfg, args)
}
