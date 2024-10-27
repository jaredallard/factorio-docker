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
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/jaredallard/factorio-docker/internal/config"
)

// runVanilla runs the Factorio server as normal.
func runVanilla(ctx context.Context, cfg *config.Config, args []string) error {
	//nolint:gosec // Why: We're creating the arguments above.
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Dir = cfg.InstallPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Start the process in a new process group so we can kill it and all
	// of its children reliably. This also detaches ^C (sent to us) from
	// killing the child process, instead allowing our context cancel to
	// do that.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return cmd.Process.Signal(syscall.SIGINT)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start '%s': %w", cmd.String(), err)
	}

	if err := cmd.Wait(); err != nil {
		// Only report errors if the context wasn't canceled.
		if ctx.Err() == nil {
			return fmt.Errorf("failed to run '%s': %w", cmd.String(), err)
		}
	}

	return nil
}
