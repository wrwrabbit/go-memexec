package memexec

import (
	"os"
	"os/exec"
)

// Exec is an in-memory executable code unit.
type Exec struct {
	executor
}

func NewDirect(b []byte, filename string) (*Exec, error) {
	f, err := DirectTempFile("", filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = f.Close()
			_ = os.Remove(f.Name())
		}
	}()

	return PrepareTemp(b, f)
}

// New creates new memory execution object that can be
// used for executing commands on a memory based binary.
func New(b []byte) (*Exec, error) {
	f, err := TempFile("", tempPattern)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = f.Close()
			_ = os.Remove(f.Name())
		}
	}()

	return PrepareTemp(b, f)
}

func PrepareTemp(b []byte, f *os.File) (*Exec, error) {
	// we need only read and execution privileges
	// ioutil.TempFile creates files with 0600 perms
	if err := os.Chmod(f.Name(), 0600); err != nil {
		return nil, err
	}
	if _, err := f.Write(b); err != nil {
		return nil, err
	}

	var exe Exec
	if err := exe.prepare(f); err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	return &exe, nil
}

// Command is an equivalent of `exec.Command`,
// except that the path to the executable is be omitted.
func (m *Exec) Command(arg ...string) *exec.Cmd {
	return exec.Command(m.path(), arg...)
}

// Close closes Exec object.
//
// Any further command will fail, it's client's responsibility
// to control the flow by using synchronization algorithms.
func (m *Exec) Close() error {
	return m.close()
}
