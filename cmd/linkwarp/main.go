package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/warpfork/linkwarp"
)

func main() {
	var searchRoot, binRoot string
	switch len(os.Args) {
	case 1:
		self, err := os.Executable()
		if err != nil {
			fmt.Printf("linkwarp: cannot find self: %s\n", err)
			os.Exit(26)
		}
		os.Chdir(filepath.Dir(self))
		// Go up three, because we want to exit "bin/", exit (presumably) "linkwarp*/", exit (presumably) "apps/", and then enter "apps/" and "bin/" again explicitly from there.
		// This results in linkwarp DTRT'ing if installed in the canonical sort of location warpsys conventions would have you put it in, and in all links still being relative.
		searchRoot = "../../../apps/"
		binRoot = "../../../bin/"
	case 2:
		searchRoot = filepath.Join(os.Args[1], "apps")
		binRoot = filepath.Join(os.Args[1], "bin")
	case 3:
		searchRoot = os.Args[1]
		binRoot = os.Args[2]
	default:
		fmt.Printf("linkwarp: incorrect usage, acceptable usage is 0, 1, or 2 args\n")
		os.Exit(27)
	}

	searchCfg := linkwarp.BinSearchCfg{
		SearchRoot: searchRoot,
	}
	synthCfg := linkwarp.BinSynthesisCfg{
		BinRoot: binRoot,
		AppRoot: searchCfg.SearchRoot,
	}
	if err := searchCfg.StartSearch(context.Background(), synthCfg.UpdateLinks); err != nil {
		fmt.Printf("linkwarp: error: %s\n", err)
		os.Exit(25)
	}
}
