package runc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/urfave/cli/v2"
)

var RunCommand = &cli.Command{
	Name: "run",
	Usage: `启动一个隔离的容器
			gocker run -it [command]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "it",
			Usage: "是否启用命令行交互模式",
		},
		&cli.StringFlag{
			Name:  "rootfs",
			Usage: "容器根目录",
		},
	},
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return errors.New("参数不全，请检查！")
		}
		read, write, err := os.Pipe()
		if err != nil {
			return err
		}

		tty := context.Bool("it")
		rootfs := context.String("rootfs")

		cmd := exec.Command("/proc/self/exe", "init")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWNS |
				syscall.CLONE_NEWUTS |
				syscall.CLONE_NEWIPC |
				syscall.CLONE_NEWPID |
				syscall.CLONE_NEWNET,
		}
		if tty {
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		cmd.ExtraFiles = []*os.File{read}
		cmd.Dir = rootfs

		if err := cmd.Start(); err != nil {
			log.Println("command start error", err)
			return err
		}
		write.WriteString(strings.Join(context.Args().Slice(), " "))
		write.Close()

		cmd.Wait()
		return nil
	},
}

var InitCommand = &cli.Command{
	Name:  "init",
	Usage: "初始化容器进程，请勿直接调用！",
	Action: func(context *cli.Context) error {
		pwd, err := os.Getwd()
		if err != nil {
			log.Printf("Get current path error %v", err)
			return err
		}

		log.Println("Current path is ", pwd)
		cmdArray := readCommandArray()
		if cmdArray == nil || len(cmdArray) == 0 {
			return fmt.Errorf("Command is empty")
		}

		log.Println("CmdArray is ", cmdArray)

		err = pivotRoot(pwd)
		if err != nil {
			log.Printf("pivotRoot error %v", err)
			return err
		}

		//mount proc
		defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
		syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

		// 配置hostname
		if err := syscall.Sethostname([]byte("lagoudocker")); err != nil {
			fmt.Printf("Error setting hostname - %s\n", err)
			return err
		}

		path, err := exec.LookPath(cmdArray[0])
		if err != nil {
			log.Printf("Exec loop path error %v", err)
			return err
		}

		// export PATH=$PATH:/bin
		if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
			log.Println(err.Error())
		}

		return nil
	},
}

func pivotRoot(root string) error {
	// 确保新 root 和老 root 不在同一目录
	// MS_BIND：执行bind挂载，使文件或者子目录树在文件系统内的另一个点上可视。
	// MS_REC： 创建递归绑定挂载，递归更改传播类型
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}

	// 创建 .pivot_root 文件夹，用于存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	// 调用 Golang 封装的 PivotRoot
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	// 修改工作目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// 卸载 .pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	// 删除临时文件夹 .pivot_root
	return os.Remove(pivotDir)
}

func readCommandArray() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Printf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
