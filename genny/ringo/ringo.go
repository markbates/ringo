package ringo

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	if err := g.Box(packr.New("github.com/markbates/ringo/genny/ringo", "../ringo")); err != nil {
		return g, errors.WithStack(err)
	}

	if err := g.Box(packr.New("github.com/markbates/ringo/genny/ringo/templates", "../ringo/templates")); err != nil {
		return g, errors.WithStack(err)
	}

	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("-name-", strings.ToLower(opts.Name)))
	g.Transformer(gogen.TemplateTransformer(opts, nil))
	g.Transformer(up())

	g.Command(exec.Command("go", "version"))

	g.RunFn(buffalo())
	g.RunFn(goEnv())

	return g, nil
}

func goEnv() genny.RunFn {
	return func(r *genny.Runner) error {
		bb := &bytes.Buffer{}
		c := exec.Command("go", "env")
		c.Stdout = bb
		c.Stderr = bb
		if err := r.Exec(c); err != nil {
			return err
		}
		f := genny.NewFile("-name-.env.log.up", bb)
		return r.File(f)
	}
}

func buffalo() genny.RunFn {
	return func(r *genny.Runner) error {
		req, err := http.NewRequest("GET", "https://www.gobuffalo.io/en", nil)
		if err != nil {
			return err
		}
		res, err := r.Request(req)
		if err != nil {
			return err
		}
		if res == nil {
			r.Logger.Debug("~~no response~~")
			return nil
		}
		// defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil && err != context.Canceled {
			return err
		}
		f := genny.NewFileB("index.html", b)
		return r.File(f)
	}
}

func up() genny.Transformer {
	t := genny.NewTransformer(".up", func(f genny.File) (genny.File, error) {
		return genny.NewFileS(f.Name(), strings.ToUpper(f.String())), nil
	})
	t.StripExt = true
	return t
}
