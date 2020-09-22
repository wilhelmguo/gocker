package cgroups

import (
	"fmt"
	"os"
	"path"
)

func getCgroupPath(subsystem string, autoCreate bool) (string, error) {
	cgroupRootSubsystem := path.Join(cgroupsRoot, subsystem)
	_, err := os.Stat(path.Join(cgroupRootSubsystem, gockerCgroupPath))
	if err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRootSubsystem, gockerCgroupPath), 0755); err == nil {
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRootSubsystem, gockerCgroupPath), nil
	}
	return "", fmt.Errorf("cgroup path error %v", err)
}
