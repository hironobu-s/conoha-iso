package command

import (
	"io"
	"io/ioutil"
)

type Command struct {
}

func (cmd *Command) extractApiErrorMessage(r io.Reader) string {
	errjson, err := ioutil.ReadAll(r)
	if err != nil {
		return err.Error()
	}
	return string(errjson)
}
