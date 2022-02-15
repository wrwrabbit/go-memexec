//+build linux

package memexec

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type executor struct {
	f *os.File
}

// on linux we can keep a read only fd of the temp file and remove it,
// kernel buffers its content in memory until all fds are closed.
func (e *executor) prepare(t *os.File) error {
	f, err := os.OpenFile(t.Name(), os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = f.Close()
		}
	}()

	// check if /proc is mounted
	path := fmt.Sprintf("/proc/self/fd/%d", int(f.Fd()))
	if _, err := os.Lstat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s dosn't exist, probably /proc is not mounted", path)
		}
		return err
	}
	if err = os.Remove(t.Name()); err != nil {
		return err
	}

	e.f = f
	return nil
}

func (e *executor) path() string {
	return fmt.Sprintf("/proc/self/fd/%d", int(e.f.Fd()))
}

func (e *executor) close() error {
	return e.f.Close()
}

func makeExecutive(targetExePath string) {
	showError(executeOnce(fmt.Sprintf("chmod a+x %s", targetExePath)))
}

func executeOnce(command string) error {
	shellInterpreter := "/bin/sh"

	cmd := exec.Command(shellInterpreter)
	pipeIn, _ := cmd.StdinPipe()

	pipeIn.Write([]byte(command + " & exit\n"))
	_, err := cmd.CombinedOutput()
	return err
}

func showError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
