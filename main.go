package main

import (
	"context"
	"flag"
	"log"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/logger"
	"github.com/markbates/ringo/genny/ringo"
)

var dryRun bool
var verbose bool

func main() {
	flag.BoolVar(&dryRun, "dry-run", false, "runs the generator dry")
	flag.BoolVar(&verbose, "verbose", false, "run with verbose output")
	flag.Parse()

	ctx := context.Background()

	run := genny.WetRunner(ctx)
	if dryRun {
		run = genny.DryRunner(ctx)
	}

	if verbose {
		run.Logger = logger.New(logger.DebugLevel)
	}

	if err := run.WithNew(ringo.New(&ringo.Options{Name: "Ringo"})); err != nil {
		log.Fatal(err)
	}

	if err := run.Run(); err != nil {
		log.Fatal(err)
	}

}
