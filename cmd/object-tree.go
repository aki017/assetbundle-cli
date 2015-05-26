package cmd

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/aki017/assetbundle"
	"github.com/aki017/assetbundle/body"
	"github.com/codegangsta/cli"
)

// ObjectTree is TypeTree command
func ObjectTree(c *cli.Context) {
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
						showObjectTree(body, v)
					} else {
						logrus.Fatal("ClassID1 != ClassID2")
					}
				}
			}
		}
	}
}

func showObjectTree(b *body.Body, o *body.Object) {
	data := b.GetBody()[b.Header.DataOffset+o.Offset : b.Header.DataOffset+o.Offset+o.Length]
	fmt.Printf(
		"      %d: {offset: 0x%08x, length: 0x%08x, Class: %16s}\n",
		o.ID,
		o.Offset,
		o.Length,
		o.ClassID1.String(),
	)
	objecttree := b.TypeTree.Get(o.ClassID1).Parse(data)
	//s, _ := json.MarshalIndent(objecttree, "", "  ")

	s := objecttree.String()
	fmt.Printf("%s\n", s)
}

// CmdObjectTree is command
var CmdObjectTree = cli.Command{
	Name:      "object-tree",
	ShortName: "ot",
	Usage:     "objecttreecommand",
	Action:    ObjectTree,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "debug",
		},
	},
}
