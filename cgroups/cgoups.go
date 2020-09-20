package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

const cpuPath = "/sys/fs/cgroup/cpu/gocker"
const memoryPath = "/sys/fs/cgroup/memory/gocker"

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
	if CPU != 0 {
		err := ioutil.WriteFile(path.Join(cpuPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
		if err != nil {
			return fmt.Errorf("set cgroup cpu fail %v", err)
		}
	}
	if Memory != 0 {
		err := ioutil.WriteFile(path.Join(memoryPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
		if err != nil {
			return fmt.Errorf("set cgroup memory fail %v", err)
		}
	}
	return nil
}

// 释放cgroup
func (c *Cgroups) Destroy() error {
	if CPU != 0 {
		return os.RemoveAll(cpuPath)
	}
	if Memory != 0 {
		return os.RemoveAll(memoryPath)
	}
	return nil
}

func (c *Cgroups) SetCPULimit(cpu int) error {
	if err := ioutil.WriteFile(path.Join(cpuPath, "cpu.cfs_quota_us"), []byte(strconv.Itoa(cpu*100000)), 0644); err != nil {
		return fmt.Errorf("set cpu limit fail %v", err)
	}
	return nil
}

func (c *Cgroups) SetMemoryLimit(memory int) error {
	if err := ioutil.WriteFile(path.Join(cpuPath, "memory.limit_in_bytes"), []byte(strconv.Itoa(memory*1024*1024)), 0644); err != nil {
		return fmt.Errorf("set memory limit fail %v", err)
	}
	return nil
}
