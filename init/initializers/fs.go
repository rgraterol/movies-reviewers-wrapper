package initializers

import (
	"os"
)

type fsGetter interface {
	getwd() (dir string, err error)
}

type fsGetterFunc func() (dir string, err error)

func (f fsGetterFunc) getwd() (dir string, err error) {
	return f()
}

var osFSGetter = fsGetterFunc(os.Getwd)

