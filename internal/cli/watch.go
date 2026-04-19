package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/nicholasgasior/envchain-go/internal/store"
)

// CmdWatch runs a command with the project's env vars, reloading on SIGHUP.
// Usage: envchain watch <project> [--] <cmd> [args...]
func CmdWatch(st *store.Store, project, passphrase string, args []string, out io.Writer) error {
	if project == "" || len(args) == 0 {
		return fmt.Errorf("usage: envchain watch <project> [--] <cmd> [args...]")
	}

	run := func() (*exec.Cmd, error) {
		set, err := st.Load(project, passphrase)
		if err != nil {
			return nil, fmt.Errorf("load project: %w", err)
		}
		env := os.Environ()
		for _, k := range set.Keys() {
			v, _ := set.Get(k)
			env = append(env, k+"="+v)
		}
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Env = env
		cmd.Stdout = out
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("start command: %w", err)
		}
		fmt.Fprintf(out, "[watch] started pid %d at %s\n", cmd.Process.Pid, time.Now().Format(time.RFC3339))
		return cmd, nil
	}

	cmd, err := run()
	if err != nil {
		return err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	for {
		done := make(chan error, 1)
		go func() { done <- cmd.Wait() }()

		select {
		case sig := <-sigs:
			switch sig {
			case syscall.SIGHUP:
				fmt.Fprintf(out, "[watch] reloading...\n")
				_ = cmd.Process.Kill()
				<-done
				cmd, err = run()
				if err != nil {
					return err
				}
			default:
				_ = cmd.Process.Kill()
				<-done
				fmt.Fprintf(out, "[watch] stopped\n")
				return nil
			}
		case exitErr := <-done:
			if exitErr != nil {
				fmt.Fprintf(out, "[watch] process exited: %v\n", exitErr)
			}
			return nil
		}
	}
}
