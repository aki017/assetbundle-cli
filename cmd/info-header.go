package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/aki017/assetbundle"
	"github.com/aki017/assetbundle/header"
	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
)

func parse(path string) *header.Header {
	fp, err := os.Open(path)
	if err != nil {
		log.Panic(err)
	}
	defer fp.Close()
	b, _ := ioutil.ReadAll(fp)
	return assetbundle.DecodeHeader(b)
}

// InfoHeader is file info command
func InfoHeader(c *cli.Context) {
	bundles := make(map[string]*header.Header)
	for _, path := range c.Args() {
		bundles[path] = parse(path)
	}

	switch c.GlobalString("format") {
	case "json":
		j, _ := json.Marshal(bundles)
		fmt.Println(string(j))
	case "prettyjson":
		j, _ := json.MarshalIndent(bundles, "", "  ")
		fmt.Println(string(j))
	default:
		for path, header := range bundles {
			fmt.Println("#", path)
			fmt.Println("FileType:", header.FileType.String())
			fmt.Println("Format:", header.Format.String())
			fmt.Println("PlayerVersion:", header.PlayerVersion)
			fmt.Println("EngineVersion:", header.EngineVersion)
			fmt.Println("FileSize:", humanize.Bytes(uint64(header.FileSize)))
			fmt.Println("DataOffset:", strconv.FormatUint(uint64(header.DataOffset), 10))
			fmt.Println("Unknown1:", strconv.FormatUint(uint64(header.Unknown1), 10))
			fmt.Println()
		}
	}
}

// CmdInfoHeader is command
var CmdInfoHeader = cli.Command{
	Name:      "info-header",
	ShortName: "ih",
	Usage:     "infoheadercommand",
	Action:    InfoHeader,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "debug",
		},
	},
}
