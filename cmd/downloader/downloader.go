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

// Package main implements the downloader CLI. This CLI downloads
// Factorio and installs it.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jaredallard/factorio-docker/internal/downloader"
	"github.com/spf13/cobra"
)

// rootCmd is the root command for the downloader CLI.
var rootCmd = &cobra.Command{
	Use:   "downloader <output-dir>",
	Short: "Downloads a headless Factorio server to a directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return entrypoint(cmd, args[0])
	},
}

// entrypoint is the entrypoint for the downloader CLI.
func entrypoint(cmd *cobra.Command, outputDir string) error {
	version, _ := cmd.Flags().GetString("version")     //nolint:errcheck
	sha256sum, _ := cmd.Flags().GetString("sha256sum") //nolint:errcheck

	// If it doesn't exist, create the output directory.
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0o755); err != nil {
			return err
		}
	}

	// Ensure the output directory is empty.
	if files, err := os.ReadDir(outputDir); err == nil && len(files) > 0 {
		return fmt.Errorf("output directory %q is not empty", outputDir)
	}

	return downloader.Download(cmd.Context(), version, sha256sum, outputDir)
}

// main runs sets up and runs Cobra.
func main() {
	rootCmd.PersistentFlags().String("version", "stable",
		"The version of Factorio to download. Can be 'stable' or 'experimental' for latest release in that channel.")
	rootCmd.PersistentFlags().String("sha256sum", "",
		"The SHA256 sum of the Factorio download, if set it will be used instead of fetching it.")

	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
