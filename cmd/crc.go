package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/aki017/assetbundle"
	"github.com/codegangsta/cli"
)

// CRC is file info command
func CRC(c *cli.Context) {
	crcs := map[string]uint32{}

	var wg sync.WaitGroup
	for _, path := range c.Args() {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			ab := assetbundle.DecodeFile(p)
			if len(ab.Bodies) != 1 {
				log.Fatal("Not AssetBundle")
			}
			crcs[p] = ab.Bodies[0].CRC()
		}(path)
	}
	wg.Wait()

	switch c.GlobalString("format") {
	case "json":
		j, _ := json.Marshal(crcs)
		fmt.Println(string(j))
	case "prettyjson":
		j, _ := json.MarshalIndent(crcs, "", "  ")
		fmt.Println(string(j))
	default:
		for k, v := range crcs {
			fmt.Println(k, v)
		}
	}
}

// CmdCRC is command
var CmdCRC = cli.Command{
	Name:      "crc",
	ShortName: "c",
	Usage:     "crccommand",
	Action:    CRC,
}
