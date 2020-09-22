package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

const gockerCgroupPath = "gocker"
const cgroupsRoot = "/sys/fs/cgroup"

type Cgroups struct {
	// 单位 核
	CPU int
	// 单位 兆
	Memory int
}

func NewCgroups() *Cgroups {
	return &Cgroups{}
}

func (c *Cgroups) Apply(pid int) error {
	if c.CPU != 0 {
		cpuCgroupPath, err := getCgroupPath("cpu", true)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path.Join(cpuCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
		if err != nil {
			return fmt.Errorf("set cgroup cpu fail %v", err)
		}
	}
	if c.Memory != 0 {
		memoryCgroupPath, err := getCgroupPath("memory", true)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path.Join(memoryCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
		if err != nil {
			return fmt.Errorf("set cgroup memory fail %v", err)
		}
	}
	return nil
}

// 释放cgroup
func (c *Cgroups) Destroy() error {
	if c.CPU != 0 {
		cpuCgroupPath, err := getCgroupPath("cpu", false)
		if err != nil {
			return err
		}
		return os.RemoveAll(cpuCgroupPath)
	}
	if c.Memory != 0 {
		memoryCgroupPath, err := getCgroupPath("memory", false)
		if err != nil {
			return err
		}
		return os.RemoveAll(memoryCgroupPath)
	}
	return nil
}

func (c *Cgroups) SetCPULimit(cpu int) error {
	cpuCgroupPath, err := getCgroupPath("cpu", true)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(cpuCgroupPath, "cpu.cfs_quota_us"), []byte(strconv.Itoa(cpu*100000)), 0644); err != nil {
		return fmt.Errorf("set cpu limit fail %v", err)
	}
	return nil
}

func (c *Cgroups) SetMemoryLimit(memory int) error {
	memoryCgroupPath, err := getCgroupPath("memory", true)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(memoryCgroupPath, "memory.limit_in_bytes"), []byte(strconv.Itoa(memory*1024*1024)), 0644); err != nil {
		return fmt.Errorf("set memory limit fail %v", err)
	}
	return nil
}
