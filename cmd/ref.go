package cmd

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aki017/assetbundle"
	"github.com/codegangsta/cli"
)

// Ref is TypeTree command
func Ref(c *cli.Context) {
	bundles := make(map[string]*assetbundle.AssetBundle)

	var wg sync.WaitGroup
	for _, path := range c.Args() {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			bundles[p] = assetbundle.DecodeFile(p)
		}(path)
	}
	wg.Wait()

	switch c.GlobalString("format") {
	case "json":
		j, _ := json.Marshal(bundles)
		fmt.Println(string(j))
	case "prettyjson":
		j, _ := json.MarshalIndent(bundles, "", "  ")
		fmt.Println(string(j))
	default:
		for path, bundle := range bundles {
			fmt.Println("#", path)
			for _, body := range bundle.Bodies {
				fmt.Printf("%s:\n", body.Name)
				fmt.Println("  Refs:")
				fmt.Println("    Size:", len(body.AssetRefs.List))
				fmt.Println("    List:")
				for _, v := range body.AssetRefs.List {
					fmt.Printf(
						"      %s: {type: %d, `%s` `%s`}\n",
						v.GUID.String(),
						v.Type,
						v.AssetPath,
						v.FilePath)
				}
			}
		}
	}
}

// CmdRef is command
var CmdRef = cli.Command{
	Name:      "ref",
	ShortName: "r",
	Usage:     "refcommand",
	Action:    Ref,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "debug",
		},
	},
}
