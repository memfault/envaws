package runner

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func RunCmd(cmd *exec.Cmd) {
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()

	if err != nil {
		log.Fatal(err)
	}

	waitCh := make(chan error, 1)

	go func() {
		waitCh <- cmd.Wait()
		close(waitCh)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh)

	// loop to handle multiple signals
	for {
		select {
		case sig := <-sigCh:
			if sig == syscall.SIGCHLD {
				continue
			}
			if err := cmd.Process.Signal(sig); err != nil {
				log.Print("error sending signal", sig, err)
			}
		case err := <-waitCh:
			// Subprocess exited, get the return code if we can
			var waitStatus syscall.WaitStatus
			if exitError, ok := err.(*exec.ExitError); ok {
				waitStatus = exitError.Sys().(syscall.WaitStatus)
				os.Exit(waitStatus.ExitStatus())
			}
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
}

// Send a SIGTERM (soft-kill) to the worker
func SoftThenHardKill(cmd *exec.Cmd, timeout time.Duration) {
	if cmd != nil {
		log.Println("Sending SIGTERM...")
		cmd.Process.Signal(syscall.SIGTERM)
	}
	time.Sleep(timeout)
	if cmd != nil {
		log.Println("Sending SIGKILL...")
		cmd.Process.Signal(syscall.SIGKILL)
	}
}
