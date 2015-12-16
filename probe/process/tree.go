package process

import (
	"fmt"
)

// Tree represents all processes on the machine.
type Tree interface {
	GetParent(pid int) (int, error)
	GetChildren(pid int) ([]int, error)
}

type tree struct {
	processes map[int]Process
}

// NewTree returns a new Tree that can be polled.
func NewTree(walker Walker) (Tree, error) {
	pt := tree{processes: map[int]Process{}}
	err := walker.Walk(func(p Process) {
		pt.processes[p.PID] = p
	})

	return &pt, err
}

// GetParent returns the pid of the parent process for a given pid
func (pt *tree) GetParent(pid int) (int, error) {
	proc, ok := pt.processes[pid]
	if !ok {
		return -1, fmt.Errorf("PID %d not found", pid)
	}

	return proc.PPID, nil
}

// GetChildren
func (pt *tree) GetChildren(pid int) ([]int, error) {
	_, ok := pt.processes[pid]
	if !ok {
		return []int{}, fmt.Errorf("PID %d not found", pid)
	}

	var isChild func(id int) bool
	isChild = func(id int) bool {
		p, ok := pt.processes[id]
		if !ok || p.PPID == 0 {
			return false
		} else if p.PPID == pid {
			return true
		}
		return isChild(p.PPID)
	}

	children := []int{pid}
	for id := range pt.processes {
		if isChild(id) {
			children = append(children, id)
		}
	}
	return children, nil
}
