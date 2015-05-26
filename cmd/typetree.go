package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/aki017/assetbundle"
	"github.com/aki017/assetbundle/body"
	"github.com/codegangsta/cli"
)

// TypeTree is TypeTree command
func TypeTree(c *cli.Context) {
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
				fmt.Println("  TypeTree:")
				fmt.Println("    Version:", body.TypeTree.Version)
				fmt.Println("    UnityVersion:", body.TypeTree.UnityVersion)
				fmt.Println("    FieldSize:", len(body.TypeTree.Fields))
				fmt.Println("    Fields:")
				for _, f := range body.Objects.List {
					fmt.Println(format(body.TypeTree.Get(f.ClassID1), body))
				}
			}
		}
	}
}

func format(f *body.TypeField, b *body.Body) string {
	cs := make([]string, len(f.Children))
	for i, f2 := range f.Children {
		cs[i] = "  " + strings.Replace(format(f2, b), "\n", "\n  ", -1)
	}
	children := strings.Join(cs, "\n")
	re := fmt.Sprintf(
		"%s: %s (size: %d, index: %d, ArrayFlag: %d, %d %d)",
		f.Name,
		f.Type,
		f.Size,
		f.Index,
		f.ArrayFlag,
		f.Flags1,
		f.Flags2,
	)
	if (len(cs)) != 0 {
		re += "\n" + children
	}
	return re
}

// CmdTypeTree is command
var CmdTypeTree = cli.Command{
	Name:      "typetree",
	ShortName: "t",
	Usage:     "typetreecommand",
	Action:    TypeTree,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "debug",
		},
	},
}
