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

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	charmlog "github.com/charmbracelet/log"

	"github.com/jaredallard/factorio-docker/internal/config"
	"github.com/jaredallard/factorio-docker/internal/downloader"
	"github.com/jaredallard/factorio-docker/internal/launcher"
)

// entrypoint is the entrypoint fro the wrapper CLI.
func entrypoint() error {
	handler := charmlog.New(os.Stderr)
	log := slog.New(handler)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log.Info("Starting...", "version", cfg.Version, "server_path", cfg.InstallPath, "data_path", cfg.ServerDataPath)

	// Ensure that we're using the requested version.
	log.Info("Checking installed Factorio version")
	if err := downloader.EnsureVersion(ctx, cfg, log); err != nil {
		return err
	}

	// Launch the Factorio server.
	return launcher.Launch(ctx, log, cfg)
}

// main runs the entrypoint function. If it returns a non-nil error, it
// exits the program with a status code of 1. This is done to prevent
// `defer` from swallowing panics.
func main() {
	if err := entrypoint(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
