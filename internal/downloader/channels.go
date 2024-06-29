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
	"fmt"

	"github.com/jaredallard/factorio-docker/internal/factorio"
)

// Channel is a Factorio release channel.
type Channel string

const (
	// ChannelStable is the stable channel.
	ChannelStable Channel = "stable"

	// ChannelExperimental is the experimental channel.
	ChannelExperimental Channel = "experimental"
)

// String returns the string representation of the channel.
func (c Channel) String() string {
	return string(c)
}

// VersionIsChannel returns true if the specified version is a channel
func VersionIsChannel(version string) bool {
	return version == ChannelStable.String() || version == ChannelExperimental.String()
}

// GetVersionForChannel returns the latest version of the specified
// channel.
func GetVersionForChannel(channel Channel) (string, error) {
	// Special case, resolve stable/experimental to the latest version of
	// their respective channels.
	rels, err := factorio.GetLatestReleases()
	if err != nil {
		return "", err
	}

	var version string
	switch channel {
	case ChannelStable:
		version = rels.Stable
	case ChannelExperimental:
		version = rels.Experimental
	default:
		return "", fmt.Errorf("unknown channel %s", channel)
	}

	return version, nil
}
