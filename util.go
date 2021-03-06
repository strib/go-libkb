package libkb

import (
	"os"
	"path"
)

func ErrToOk(err error) string {
	if err == nil {
		return "ok"
	} else {
		return "ERROR"
	}
}

// exists returns whether the given file or directory exists or not
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func MakeParentDirs(filename string) error {

	dir, _ := path.Split(filename)
	exists, err := FileExists(dir)
	if err != nil {
		G.Log.Error("Can't see if parent dir %s exists", dir)
		return err
	}

	if !exists {
		err = os.MkdirAll(dir, PERM_DIR)
		if err != nil {
			G.Log.Error("Can't make parent dir %s", dir)
			return err
		} else {
			G.Log.Info("Created parent directory %s", dir)
		}
	}
	return nil
}
