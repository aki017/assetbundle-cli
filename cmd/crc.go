package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aki017/assetbundle"
	"github.com/codegangsta/cli"
)

// CRC is file info command
func CRC(c *cli.Context) {
	crcs := map[string]uint32{}
	m := new(sync.Mutex)

	var wg sync.WaitGroup
	var files []string
	if c.Bool("stdin") {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			files = append(files, scanner.Text())
		}
	} else {
		files = c.Args()
	}
	for _, path := range files {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			ab := assetbundle.DecodeFile(p)
			if len(ab.Bodies) != 1 {
				log.Fatal("Not AssetBundle")
			}
			m.Lock()
			crcs[p] = ab.Bodies[0].CRC()
			m.Unlock()
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
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "stdin",
			Usage: "read files from stdin",
		},
	},
}
