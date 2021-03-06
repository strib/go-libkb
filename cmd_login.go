package libkb

import ()

type CmdLogin struct{}

func (v CmdLogin) Run() error {

	err := G.LoginState.Login()
	if err != nil {
		return err
	}
	return nil
}

func (v CmdLogin) UseConfig() bool   { return true }
func (v CmdLogin) UseKeyring() bool  { return true }
func (v CmdLogin) UseAPI() bool      { return true }
func (v CmdLogin) UseTerminal() bool { return true }
