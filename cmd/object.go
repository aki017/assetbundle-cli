package cmd

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aki017/assetbundle"
	"github.com/codegangsta/cli"
)

// Object is TypeTree command
func Object(c *cli.Context) {
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
				fmt.Println("  Objects:")
				fmt.Println("    Size:", len(body.Objects.List))
				fmt.Println("    List:")
				for _, v := range body.Objects.List {
					if v.ClassID1 == v.ClassID2 {
						fmt.Printf(
							"      %d: {offset: 0x%08x, length: 0x%08x, Class: %16s}\n",
							v.ID,
							v.Offset,
							v.Length,
							v.ClassID1.String(),
						)
					} else {
						fmt.Printf(
							"      %d: {offset: 0x%08x, length: 0x%08x, ClassID1: %8d, ClassID2: %8d}\n",
							v.ID,
							v.Offset,
							v.Length,
							v.ClassID1,
							v.ClassID2,
						)
					}
				}
			}
		}
	}
}

// CmdObject is command
var CmdObject = cli.Command{
	Name:      "object",
	ShortName: "o",
	Usage:     "objectcommand",
	Action:    Object,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "debug",
		},
	},
}
