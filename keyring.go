package libkb

import (
	"code.google.com/p/go.crypto/openpgp"
	"fmt"
	"os"
)

type KeyringFile struct {
	filename string
	Entities openpgp.EntityList
	isPublic bool
}

type Keyrings struct {
	Public []KeyringFile
	Secret []KeyringFile
}

func (k Keyrings) MakeKeyrings(out *[]KeyringFile, filenames []string, isPublic bool) {
	*out = make([]KeyringFile, len(filenames))
	for i, filename := range filenames {
		(*out)[i] = KeyringFile{filename, openpgp.EntityList{}, isPublic}
	}
}

func NewKeyrings(e Env) *Keyrings {
	ret := &Keyrings{}
	ret.MakeKeyrings(&ret.Public, e.GetPublicKeyrings(), true)
	ret.MakeKeyrings(&ret.Secret, e.GetSecretKeyrings(), false)
	return ret
}

func (k *Keyrings) Load() (err error) {
	G.Log.Debug("+ Loading keyrings")
	err = k.LoadKeyrings(k.Public)
	if err == nil {
		k.LoadKeyrings(k.Secret)
	}
	G.Log.Debug("- Loaded keyrings")
	return err
}

func (k Keyrings) LoadKeyrings(v []KeyringFile) (err error) {
	for _, k := range v {
		if err = k.Load(); err != nil {
			return err
		}
	}
	return nil
}

func (k *KeyringFile) Load() error {
	G.Log.Debug(fmt.Sprintf("+ Loading PGP Keyring %s", k.filename))
	file, err := os.Open(k.filename)
	if os.IsNotExist(err) {
		G.Log.Warning(fmt.Sprintf("No PGP Keyring found at %s", k.filename))
		err = nil
	} else if err != nil {
		G.Log.Error(fmt.Sprintf("Cannot open keyring %s: %s\n", err.Error()))
		return err
	}
	if file != nil {
		k.Entities, err = openpgp.ReadKeyRing(file)
		if err != nil {
			G.Log.Error(fmt.Sprintf("Cannot parse keyring %s: %s\n", err.Error()))
			return err
		}
	}
	G.Log.Debug(fmt.Sprintf("- Successfully loaded PGP Keyring"))
	return nil
}

func (k KeyringFile) writeTo(file *os.File) error {
	for _, e := range k.Entities {
		if err := e.Serialize(file); err != nil {
			return err
		}
	}
	return nil
}

func (k KeyringFile) Save() error {
	G.Log.Debug(fmt.Sprintf("+ Writing to PGP keyring %s", k.filename))
	tmpfn, tmp, err := TempFile(k.filename, PERM_FILE)
	G.Log.Debug(fmt.Sprintf("| Temporary file generated: %s", tmpfn))
	if err != nil {
		return err
	}

	err = k.writeTo(tmp)
	if err == nil {
		err = tmp.Close()
		if err != nil {
			err = os.Rename(tmpfn, k.filename)
		} else {
			G.Log.Error(fmt.Sprintf("Error closing temporary file %s: %s", tmp, err.Error()))
			os.Remove(tmpfn)
		}
	} else {
		G.Log.Error(fmt.Sprintf("Error writing temporary keyring %s: %s", tmp, err.Error()))
		tmp.Close()
		os.Remove(tmpfn)
	}
	G.Log.Debug(fmt.Sprintf("- Wrote to PGP keyring %s -> %s", k.filename, ErrToOk(err)))
	return err
}
