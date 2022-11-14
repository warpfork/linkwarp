package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/warptools/linkwarp"

	"github.com/jawher/mow.cli"
)

type ExitCode int

const (
	ExitCodeOkay ExitCode = iota
	ExitCodeUnknownError
	ExitCodePanic
	ExitCodeError    ExitCode = 25
	ExitCodeMetanoia ExitCode = 26
	ExitCodeUsage    ExitCode = 27
)

// Go up three, because we want to exit "bin/", exit (presumably) "linkwarp*/", exit (presumably) "apps/", and then enter "apps/" and "bin/" again explicitly from there.
// This results in linkwarp DTRT'ing if installed in the canonical sort of location warpsys conventions would have you put it in, and in all links still being relative.
const (
	DefaultAppsPath = "../../../apps/"
	DefaultBinPath  = "../../../bin/"
)

// 1 arg form: linkwarp
// 2 arg form: linkwarp [root]
// 3 arg form: linkewarp [search root] [bin root]
func main() {
	ctx := &ActionContext{
		Context: context.Background(),
	}
	app := cli.App("linkwarp", "does things!\nWith no arguments will start at the path of the linkwarp executable")
	app.Spec = "[-v] ( [ROOT] | ( SEARCHROOT BINROOT ) )"

	app.BoolOptPtr(&ctx.Config.Verbose, "v verbose", false, "Enable logging output on stderr")
	app.StringArgPtr(&ctx.Config.Root, "ROOT", "", "Single argument replaces both SEARCHROOT and BINROOT")
	app.StringPtr(&ctx.Config.Search.SearchRoot, cli.StringArg{
		Name:      "SEARCHROOT",
		Value:     DefaultAppsPath,
		Desc:      "the path to search for binaries",
		SetByUser: &ctx.Config.SearchSetByUser,
	})
	app.StringArgPtr(&ctx.Config.Synthesis.BinRoot, "BINROOT", DefaultBinPath, "the path to place links")
	app.Before = Action(ctx, logConfig)
	app.Action = Action(ctx, run)
	app.After = Action(ctx, exitHandler)
	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
		os.Exit(int(ExitCodeUsage))
	}
	os.Exit(int(ExitCodeOkay))
}

type Config struct {
	Root            string
	SearchSetByUser bool
	Search          linkwarp.BinSearchCfg
	Synthesis       linkwarp.BinSynthesisCfg
	Verbose         bool
}

type ActionContext struct {
	context.Context
	Config
	Error error
}

type ActionFunc func(*ActionContext) error

func Action(c *ActionContext, funcs ...ActionFunc) func() {
	if c.Context == nil {
		c.Context = context.Background()
	}
	if c.Error != nil {
		return func() {}
	}
	return func() {
		for _, f := range funcs {
			err := f(c)
			if err != nil {
				c.Error = err
				break
			}
		}
	}
}

func logConfig(c *ActionContext) error {
	if c.Config.Verbose {
		log.SetFlags(0)
		return nil
	}
	log.SetOutput(io.Discard)
	return nil
}

func run(c *ActionContext) error {
	ctx := c.Context
	if c.Config.Root == "" && !c.Config.SearchSetByUser {
		log.Println("Args: 0")
		self, err := os.Executable()
		if err != nil {
			log.Printf("linkwarp: cannot find self: %s\n", err)
			os.Exit(int(ExitCodeMetanoia)) // TODO: proper error handling and not exiting in the middle of an action
		}
		os.Chdir(filepath.Dir(self))
		c.Config.Search.SearchRoot = DefaultAppsPath
		c.Config.Synthesis.BinRoot = DefaultBinPath
	}
	if c.Config.Root != "" {
		log.Println("Args: 1")
		c.Config.Search.SearchRoot = c.Config.Root
		c.Config.Synthesis.BinRoot = c.Config.Root
	}
	c.Config.Synthesis.AppRoot = c.Config.Search.SearchRoot
	log.Printf(`{"approot": %q, "binroot": %q}`, c.Config.Synthesis.AppRoot, c.Config.Synthesis.BinRoot)
	return c.Config.Search.StartSearch(ctx, c.Config.Synthesis.UpdateLinks)
}

func exitHandler(c *ActionContext) error {
	if c.Error != nil {
		log.Println(c.Error)
		os.Exit(int(ExitCodeError))
	}
	os.Exit(int(ExitCodeOkay))
	return nil
}
