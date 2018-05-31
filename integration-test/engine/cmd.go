package engine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/alecthomas/template"
	"github.com/loomnetwork/loomchain/integration-test/lib"
)

var tmpl = template.New("cmd")

// func init() {
// 	tmpl = tmpl.Funcs(map[string]interface{
// 		"GetAccountAddress": GetAccountAddress,
// 	})
// }

type engineCmd struct {
	conf  lib.Config
	tests lib.TestCases
	wg    *sync.WaitGroup
	errC  chan error
}

func NewCmd(conf lib.Config, tc lib.TestCases) Engine {
	return &engineCmd{
		conf:  conf,
		tests: tc,
		wg:    &sync.WaitGroup{},
		errC:  make(chan error),
	}
}

func (e *engineCmd) Run(ctx context.Context) error {
	for _, n := range e.tests {
		// evaluate template
		t, err := tmpl.Parse(n.RunCmd)
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		err = t.Execute(buf, e.conf)
		if err != nil {
			return err
		}

		fmt.Printf("---> Buf: %s", buf.String())
		args := strings.Split(buf.String(), " ")
		if len(args) == 0 {
			return errors.New("missing command")
		}
		cmd := exec.Cmd{
			Dir:  "/Users/loomnetworklock/go/src/github.com/loomnetwork/go-loom",
			Path: args[0],
			Args: args,
		}
		out, err := cmd.Output()
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", out)
	}

	return nil
}
