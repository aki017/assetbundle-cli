package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/aki017/assetbundle"
	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
)

// Info is file info command
func Info(c *cli.Context) {
	bundles := make(map[string]assetbundle.AssetBundle)

	var wg sync.WaitGroup
	for _, path := range c.Args() {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			bundles[p] = *assetbundle.DecodeFile(p)
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
			header := bundle.Header
			fmt.Println("#", path)

			fmt.Println("Header:")
			fmt.Println("  FileType:", header.FileType.String())
			fmt.Println("  Format:", header.Format.String())
			fmt.Println("  PlayerVersion:", header.PlayerVersion)
			fmt.Println("  EngineVersion:", header.EngineVersion)
			fmt.Println("  FileSize:", humanize.Bytes(uint64(header.FileSize)))
			fmt.Println("  DataOffset:", strconv.FormatUint(uint64(header.DataOffset), 10))
			fmt.Println("  Unknown1:", strconv.FormatUint(uint64(header.Unknown1), 10))
			fmt.Println("Bodies:")
			for _, body := range bundle.Bodies {
				fmt.Printf("  %s:\n", body.Name)
				fmt.Println("    Header:")
				fmt.Println("      TreeSize:", humanize.Bytes(uint64(body.Header.TreeSize)))
				fmt.Println("      FileSize:", humanize.Bytes(uint64(body.Header.FileSize)))
				fmt.Println("      Format:", strconv.FormatUint(uint64(body.Header.Format), 10))
				fmt.Println("      DataOffset:", strconv.FormatUint(uint64(body.Header.DataOffset), 10))
				fmt.Println("      Reserved:", strconv.FormatUint(uint64(body.Header.Reserved), 10))
				fmt.Println("    TypeTree:")
				fmt.Println("      Version:", body.TypeTree.Version)
				fmt.Println("      UnityVersion:", body.TypeTree.UnityVersion)
				fmt.Println("      FieldSize:", len(body.TypeTree.Fields))
				fmt.Println("    Objects:")
				fmt.Println("      Size:", len(body.Objects.List))
				fmt.Println("    AssetRefs:")
				fmt.Println("      Size:", len(body.AssetRefs.List))
				fmt.Println()
			}
		}
	}
}

// CmdInfo is command
var CmdInfo = cli.Command{
	Name:      "info",
	ShortName: "i",
	Usage:     "infocommand",
	Action:    Info,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "debug",
		},
	},
}
