package libkb

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"strconv"
)

type CmdConfig struct {
	location bool
	reset    bool
	key      string
	value    string
}

func (v *CmdConfig) Initialize(ctx *cli.Context) error {
	v.location = ctx.Bool("location")
	v.reset = ctx.Bool("reset")
	nargs := len(ctx.Args())
	if !v.location && !v.reset &&
		nargs != 1 && nargs != 2 {
		return errors.New("incorrect config usage")
	} else {
		if nargs > 0 {
			v.key = ctx.Args()[0]
		}
		if nargs > 1 {
			v.value = ctx.Args()[1]
		}
	}
	return nil
}

func (v CmdConfig) Run() error {
	configFile := G.Env.GetConfigFilename()
	if v.location {
		G.Log.Info(fmt.Sprintf("Using config file %s", configFile))
	}

	if v.reset {
		// clear out file
		cw := G.Env.GetConfigWriter()
		cw.Reset()
		cw.Write()
		// continue on to get or set on cleared file
	}

	// TODO: validate user input?

	if v.key != "" {
		if v.value != "" {
			cw := G.Env.GetConfigWriter()
			// try to convert the value to an int, and then to a bool
			// if those don't work, use a string
			if val, e := strconv.Atoi(v.value); e == nil {
				cw.SetIntAtPath(v.key, val)
			} else if val, e := strconv.ParseBool(v.value); e == nil {
				// NOTE: this will also convert strings like 't' and 'F' to
				// a bool, which could potentially cause strange errors for
				// e.g. a user named "f"
				cw.SetBoolAtPath(v.key, val)
			} else {
				cw.SetStringAtPath(v.key, v.value)
			}
			cw.Write()
		} else {
			cr := *G.Env.GetConfig()
			// TODO: print dictionaries?
			if s, is_set := cr.GetStringAtPath(v.key); is_set {
				fmt.Printf("%s: %s\n", v.key, s)
			} else if b, is_set := cr.GetBoolAtPath(v.key); is_set {
				fmt.Printf("%s: %t\n", v.key, b)
			} else if i, is_set := cr.GetIntAtPath(v.key); is_set {
				fmt.Printf("%s: %d\n", v.key, i)
			} else {
				G.Log.Info(fmt.Sprintf("%s does not map to a value", v.key))
			}
		}
	}

	return nil
}

func (v CmdConfig) UseConfig() bool   { return true }
func (v CmdConfig) UseKeyring() bool  { return false }
func (v CmdConfig) UseAPI() bool      { return false }
func (v CmdConfig) UseTerminal() bool { return false }
