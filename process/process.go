package process

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	CHILD_ENV_KEY = "GO_CHILD_ENV_KEY"
	CHILD_ENV_VAL = "GO_CHILD_VALUE"
)

func NewProcess(parentExit bool) (*exec.Cmd, error) {
	val := os.Getenv(CHILD_ENV_KEY)
	if val == CHILD_ENV_VAL {
		// child process
		return nil, nil
	}

	childEnv := append(os.Environ(), fmt.Sprintf("%s=%s", CHILD_ENV_KEY, CHILD_ENV_VAL))
	cmd := exec.Cmd{
		Path:   os.Args[0],
		Args:   os.Args,
		Env:    childEnv,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start child process: %w", err)
	}

	if parentExit {
		os.Exit(0)
	}

	return &cmd, nil
}
