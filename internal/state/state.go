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

// Package state handles tracking the state of the Factorio server.
// Mainly concerned with which version is currently installed, so that
// the wrapper can determine if it needs to up/downgrade the server.
package state

import (
	"encoding/json"
	"fmt"
	"os"
)

// State tracks the state of a Factorio server as created by the
// wrapper.
type State struct {
	// path is the path where the state file is stored.
	path string

	// Version is the current version of Factorio installed.
	Version string `json:"version"`
}

// Open opens the state file and returns the state. If it doesn't exist,
// it will return a new state. If path is empty, it will default to
// "state.json".
func Open(path string) *State {
	if path == "" {
		path = "state.json"
	}

	var defaultState = &State{path: path}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return defaultState
	}

	f, err := os.Open(path)
	if err != nil {
		return defaultState
	}
	defer f.Close()

	var s State
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return defaultState
	}

	// Ensure the path is set.
	s.path = defaultState.path

	return &s
}

// Save saves the state to the state file.
func (s *State) Save() error {
	if s.path == "" {
		return fmt.Errorf("state path is unset")
	}

	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(s)
}
