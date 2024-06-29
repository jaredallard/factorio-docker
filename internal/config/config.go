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

// Package config stores configuration for wrapping and configuring
// Factorio.
package config

import (
	"github.com/caarlos0/env/v11"
)

// Config is the configuration for the wrapper.
// TODO(jaredallard): Generate documentation from the struct.
type Config struct {
	// Username is the factorio.com username for who owns this server.
	// This is only required for making a server public. For more
	// information, see:
	// https://wiki.factorio.com/Multiplayer#How_to_list_a_server-hosted_game_on_the_matching_server
	Username string `env:"USERNAME"`

	// Token is the factorio.com token for the server. This is only
	// required for making a server public. For more information, see:
	// https://wiki.factorio.com/Multiplayer#How_to_list_a_server-hosted_game_on_the_matching_server
	Token string `env:"TOKEN,unset"`

	// Factocord is configuration for running Factocord.
	Factocord Factocord `envPrefix:"FACTOCORD_"`

	// InstallPath is the location to install Factorio to. This is NOT
	// where the save files are stored.
	InstallPath string `env:"INSTALL_PATH" envDefault:"/opt/factorio"`

	// ServerDataPath is the location to store server data, such as
	// save files.
	ServerDataPath string `env:"SERVER_DATA_PATH" envDefault:"/data"`

	// Version is the desired version of Factorio to run. If set to
	// 'stable' or 'experimental', the latest version of that channel will
	// be used. Note that this means it will also be updated on every
	// restart. If set to a specific version, that version will be used.
	Version string `env:"VERSION" envDefault:"stable"`
}

// Factocord is the configuration for running Factocord, which provides
// Discord integration for Factorio.
type Factocord struct {
	Enabled bool `env:"ENABLED" envDefault:"false"`

	// DiscordToken is the Bot user token to use for the Discord API.
	DiscordToken string `env:"DISCORD_TOKEN,unset"`

	// DiscordChannelID is the ID of the Discord channel to send messages to.
	DiscordChannelID string `env:"DISCORD_CHANNEL_ID"`

	// DiscordUserColors determines whether to use user colors in Discord messages.
	DiscordUserColors bool `env:"DISCORD_USER_COLORS" envDefault:"true"`
}

// Load loads the configuration from the environment.
func Load() (*Config, error) {
	var cfg Config
	err := env.ParseWithOptions(&cfg, env.Options{
		// Prefix all configuration variables with FACTORIO_ to prevent
		// accidental conflicts with other environment variables.
		Prefix: "FACTORIO_",
	})
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
