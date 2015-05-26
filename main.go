package main

import (
	"os"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/aki017/assetbundle-cli/cmd"
	"github.com/codegangsta/cli"
)

func com(name string, short string, usage string, f func(*cli.Context)) cli.Command {
	return cli.Command{
		Name:      name,
		ShortName: short,
		Usage:     usage,
		Action:    f,
	}
}

func flag(name string, value string, usage string, env string) cli.Flag {
	return cli.StringFlag{
		Name:   name,
		Value:  value,
		Usage:  usage,
		EnvVar: env,
	}
}

func main() {
	log.SetOutput(os.Stdout)

	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	app := cli.NewApp()
	app.Name = "assetbundle"
	app.Usage = "Unity3d AssetBundle tools"
	app.Version = "0.0.0"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		flag("format", "text", "set format", ""),
	}

	app.Commands = []cli.Command{
		cmd.CmdInfo,
		cmd.CmdInfoHeader,
		cmd.CmdCRC,
		cmd.CmdTypeTree,
		cmd.CmdObject,
		cmd.CmdObjectTree,
		cmd.CmdRef,
		cmd.CmdConvert,
	}

	app.Run(os.Args)
	/*
		app.Run([]string{
			"assetbundle",
			"--format=plain",
			"typetree",
			"/Users/aki/Desktop/panel0.unity3d",
		})
		app.Run([]string{
			"assetbundle",
			"--format=plain",
			"object-tree",
			"/Users/aki/Desktop/panel0.unity3d",
		})
	*/
}
