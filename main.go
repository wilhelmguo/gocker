package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/wilhelmguo/gocker/runc"
)

func main() {
	app := cli.NewApp()
	app.Name = "gocker"
	app.Usage = "gocker 是 golang 编写的精简版 Docker，目的是学习 Docker 的运行原理。"

	app.Commands = []*cli.Command{
		runc.InitCommand,
		runc.RunCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
