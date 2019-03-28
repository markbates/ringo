package ringo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.NoError(gentest.CompareCommands([]string{"go version", "go env"}, res.Commands))

	r.Len(res.Files, 8)

	f, err := res.Find("ringo.go")
	r.NoError(err)
	r.Contains(f.String(), "github.com/markbates/ringo/genny/ringo")

	f, err = res.Find("starr.txt")
	r.NoError(err)
	r.Equal(f.String(), "HELLO FROM STARR STARR\n")

}

func Test_New_GoVersion(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{})
	r.NoError(err)

	run := gentest.NewRunner()
	log := run.Logger.(*gentest.Logger)
	run.With(g)

	run.ExecFn = func(c *exec.Cmd) error {
		a := strings.Join(c.Args, " ")
		if a == "go version" {
			log.Info("go1.12")
		}
		return nil
	}

	r.NoError(run.Run())
	r.Contains(log.Stream.String(), "go1.12")
}

func Test_New_GoVersion_Error(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	run.ExecFn = func(c *exec.Cmd) error {
		return fmt.Errorf("command not found %s", c.Args[0])
	}

	r.Error(run.Run())
}

func Test_up(t *testing.T) {
	r := require.New(t)

	f := genny.NewFileS("foo.go.up", "package foo")
	tr := up()
	f, err := tr.Transform(f)
	r.NoError(err)
	r.Equal("PACKAGE FOO", f.String())
	r.Equal("foo.go", f.Name())
}

func Test_buffalo(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()
	run.WithRun(buffalo())

	run.RequestFn = func(req *http.Request, c *http.Client) (*http.Response, error) {
		res := httptest.NewRecorder()
		res.WriteHeader(200)
		res.WriteString("Buffalo is awesome")
		return res.Result(), nil
	}

	r.NoError(run.Run())

	res := run.Results()

	f, err := res.Find("index.html")
	r.NoError(err)
	r.Contains(f.String(), "Buffalo")
}

func Test_New_goEnv(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	run.ExecFn = func(c *exec.Cmd) error {
		a := strings.Join(c.Args, " ")
		if a == "go env" {
			c.Stdout.Write([]byte("foo=bar"))
		}
		return nil
	}

	r.NoError(run.Run())

	res := run.Results()
	f, err := res.Find("starr.env.log")
	r.NoError(err)
	r.Contains(f.String(), "FOO=BAR")
}
