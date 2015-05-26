package cmd

import (
	"log"
	"os"
	"sync"

	"github.com/aki017/assetbundle"
	"github.com/codegangsta/cli"
)

// Convert is file convert command
func Convert(c *cli.Context) {
	var wg sync.WaitGroup
	for _, path := range c.Args() {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			writeFile(p)
		}(path)
	}
	wg.Wait()
}

func writeFile(path string) {
	ab := *assetbundle.DecodeFile(path)
	for _, b := range ab.Bodies {
		b.TypeTree.Version = 27
	}

	fp, err := os.Create(path + "." + "new")
	if err != nil {
		log.Panic(err)
	}
	defer fp.Close()

	ab.Encode(fp)
}

// CmdCRC is command
var CmdConvert = cli.Command{
	Name:      "convert",
	ShortName: "conv",
	Usage:     "convertcommand",
	Action:    Convert,
}
