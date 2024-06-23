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
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func getHeadlessArchiveNameForVersion(version string) string {
	return fmt.Sprintf("factorio_headless_x64_%s.tar.xz", version)
}

// Releases contains the latest releases of Factorio based on channels.
type Releases struct {
	Stable       string
	Experimental string
}

// releaseResponse is the response from the Factorio API.
type releaseResponse struct {
	Experimental struct {
		Alpha    string `json:"alpha"`
		Demo     string `json:"demo"`
		Headless string `json:"headless"`
	} `json:"experimental"`

	Stable struct {
		Alpha    string `json:"alpha"`
		Demo     string `json:"demo"`
		Headless string `json:"headless"`
	} `json:"stable"`
}

// GetLatestReleases gets the latest releases of Factorio.
func GetLatestReleases() (*Releases, error) {
	resp, err := http.Get("https://factorio.com/api/latest-releases")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var releases releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	return &Releases{
		Stable:       releases.Stable.Headless,
		Experimental: releases.Experimental.Headless,
	}, nil
}

// GetSHA256 gets the SHA256 hash of a Factorio version. This will only
// return the hash of the headless version.
func GetSHA256(version string) (string, error) {
	expectedFileName := getHeadlessArchiveNameForVersion(version)

	resp, err := http.Get("https://factorio.com/download/sha256sums/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		// Format: SHA256_HASH FILENAME
		line := scanner.Text()
		if strings.Contains(line, expectedFileName) {
			return strings.Fields(line)[0], nil
		}
	}

	return "", fmt.Errorf("could not find SHA256 hash for %s", version)
}
