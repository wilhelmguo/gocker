package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/urfave/cli/v2"
)

func main() {
	cmd := cli.NewApp()
	cmd.Name = "gocker"
	cmd.Usage = "gocker 是一个精简版的启动容器实现"

	cmd.Commands = []*cli.Command{
		runCommand,
	}

	if err := cmd.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var runCommand = &cli.Command{
	Name:  "run",
	Usage: `使用交互命令创建一个容器：gocker run -it [command]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "it",
			Usage: "启用交互命令",
		},
	},
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("缺少容器参数！")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("ti")
		Run(tty, cmd)
		return nil
	},
}

func Run(tty bool, command string) {
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Start(); err != nil {
		log.Println("Start commamd error.", err)
	}
	cmd.Wait()
	os.Exit(-1)
}
