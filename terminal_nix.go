// +build linux darwin

package libkb

import (
	"code.google.com/p/go.crypto/ssh/terminal"
	"fmt"
	"os"
)

type Terminal struct {
	tty          *os.File
	fd           int
	old_terminal *terminal.State
	terminal     *terminal.Terminal
}

func (t *Terminal) Init() error {
	return nil
}

func NewTerminal() *Terminal {
	return &Terminal{nil, -1, nil, nil}
}

func (t *Terminal) Startup() error {
	G.Log.Debug("Opening up /dev/tty terminal on Linux and OSX")
	file, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return err
	}
	t.tty = file
	t.fd = int(t.tty.Fd())
	t.old_terminal, err = terminal.MakeRaw(t.fd)
	if err != nil {
		return err
	}
	if t.terminal = terminal.NewTerminal(file, ""); t.terminal == nil {
		return fmt.Errorf("failed to open terminal")
	}
	return nil
}

func (t Terminal) Shutdown() error {
	if t.old_terminal != nil {
		G.Log.Debug("Restoring terminal settings")

		// XXX bug in ssh/terminal. On success, we were getting an error "errno 0";
		// so let's ignore it for now.
		terminal.Restore(t.fd, t.old_terminal)
	}
	return nil
}

func (t Terminal) PromptPassword(prompt string) (string, error) {
	return t.terminal.ReadPassword(prompt)
}

func (t Terminal) Write(s string) error {
	_, err := t.terminal.Write([]byte(s))
	return err
}

func (t Terminal) Prompt(prompt string) (string, error) {
	if len(prompt) >= 0 {
		t.Write(prompt)
	}
	return t.terminal.ReadLine()
}