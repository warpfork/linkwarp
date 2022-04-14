package main

import (
	"context"
	"fmt"
	"os"

	"github.com/warpfork/linkwarp"
)

func main() {
	searchCfg := linkwarp.BinSearchCfg{
		SearchRoot: "apps",
	}
	synthCfg := linkwarp.BinSynthesisCfg{
		BinRoot: "bin",
		AppRoot: searchCfg.SearchRoot,
	}
	if err := searchCfg.StartSearch(context.Background(), synthCfg.UpdateLinks); err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(25)
	}
}
