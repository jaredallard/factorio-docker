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
	"context"
	_ "embed" // Used w/ go:embed
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jaredallard/factorio-docker/internal/config"
)

//go:embed embed/factocord-config.json
var defaultFactocordConfig []byte

// configureFactocord configures the Factocord launcher based on the
// provided configuration.
func configureFactocord(cfg *config.Config, args []string) error {
	var diskConfig map[string]any

	// We write to install-path because FactoCord does not allow
	// configuring the location of the config file.
	configPath := filepath.Join(cfg.InstallPath, "config.json")
	if _, err := os.Stat(configPath); err != nil {
		if err := os.WriteFile(configPath, defaultFactocordConfig, 0o600); err != nil {
			return fmt.Errorf("failed to write default Factocord config: %w", err)
		}
	}

	// Read the configuration file.
	b, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read Factocord config: %w", err)
	}

	if err := json.Unmarshal(b, &diskConfig); err != nil {
		return fmt.Errorf("failed to unmarshal Factocord config: %w", err)
	}

	// Update the configuration with our settings.
	diskConfig["executable"] = args[0]
	diskConfig["launch_parameters"] = args[1:]

	if cfg.Factocord.DiscordToken != "" {
		diskConfig["discord_token"] = cfg.Factocord.DiscordToken
	}

	if cfg.Factocord.DiscordChannelID != "" {
		diskConfig["factorio_channel_id"] = cfg.Factocord.DiscordChannelID
	}

	diskConfig["discord_user_colors"] = cfg.Factocord.DiscordUserColors

	// Write the configuration back to disk.
	b, err = json.Marshal(diskConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal Factocord config: %w", err)
	}

	if err := os.WriteFile(configPath, b, 0o600); err != nil {
		return fmt.Errorf("failed to write Factocord config: %w", err)
	}

	return nil
}

// runFactocord runs Factorio through the Factocord launcher to enable
// Discord presence.
func runFactocord(ctx context.Context, cfg *config.Config, args []string) error {
	// Ensure the Factocord configuration is setup based on our
	// configuration.
	if err := configureFactocord(cfg, args); err != nil {
		return fmt.Errorf("failed to configure Factocord: %w", err)
	}

	// Run Factocord.
	newArgs := []string{
		"FactoCord-3.0",
	}
	newArgs = append(newArgs, args...)

	return runVanilla(ctx, cfg, newArgs)
}
