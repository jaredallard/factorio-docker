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

// Package main implements the wrapper CLI. It is a simple wrapper
// around configuring Factorio and running it.
package main

import "os"

// entrypoint is the entrypoint fro the wrapper CLI.
func entrypoint() error {
	return nil
}

// main runs the entrypoint function. If it returns a non-nil error, it
// exits the program with a status code of 1. This is done to prevent
// `defer` from swallowing panics.
func main() {
	if err := entrypoint(); err != nil {
		os.Exit(1)
	}
}
