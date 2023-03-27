//go:build !windows

/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/2/26 下午5:22.
 *  * Author: guojia(https://github.com/guojia99)
 */

package flock

import (
	"os"
	"sync"
	"syscall"
)

func openBlockLockFile(fileName string, start, end int64) *blockLockFile {
	return &blockLockFile{
		file:  &file{fileName: fileName},
		start: start,
		end:   end,
	}
}

type blockLockFile struct {
	m          sync.Mutex
	file       *file
	flag       int16
	start, end int64
}

func (l *blockLockFile) File() *os.File { return l.file.fh }
func (l *blockLockFile) Lock() error    { return l.lock(bSetWaiteLockOperation) }
func (l *blockLockFile) RLock() error   { return l.lock(bSetReadLockOperation) }
func (l *blockLockFile) TryLock() bool  { return l.try(bSetWaiteLockOperation) }
func (l *blockLockFile) TryRLock() bool { return l.try(bSetReadLockOperation) }
func (l *blockLockFile) Unlock() error  { return l.unlock() }

func (l *blockLockFile) lock(flag int16) (err error) {
	l.m.Lock()
	defer l.m.Unlock()

	if l.flag&flag == flag {
		return nil
	}

	if err = l.file.openFile(); err != nil {
		return err
	}

	var lock = syscall.Flock_t{
		Type:   flag,
		Whence: 0,
		Start:  l.start,
		Len:    l.end - l.start,
		Pid:    int32(os.Getpid()),
	}

	if err = syscall.FcntlFlock(l.file.fh.Fd(), blockLockSet, &lock); err != nil {
		return err
	}
	l.flag ^= flag
	return nil
}

func (l *blockLockFile) unlock() error {
	l.m.Lock()
	defer l.m.Unlock()

	if l.flag == 0 || l.file.fh == nil {
		return nil
	}
	var lock = syscall.Flock_t{
		Type:   bSetUnLockOperation,
		Whence: 0,
		Start:  l.start,
		Len:    l.end - l.start,
		Pid:    int32(os.Getpid()),
	}
	if err := syscall.FcntlFlock(l.file.fh.Fd(), blockLockSet, &lock); err != nil {
		return err
	}
	l.file.closeFile()
	l.flag = 0
	return nil
}

func (l *blockLockFile) try(flag int16) bool {
	l.m.Lock()
	defer l.m.Unlock()

	if l.flag&flag == flag {
		return true
	}
	if err := l.file.openFile(); err != nil {
		return false
	}

	var lock = syscall.Flock_t{
		Type:   flag,
		Whence: 0,
		Start:  l.start,
		Len:    l.end - l.start,
		Pid:    int32(os.Getpid()),
	}

	if err := syscall.FcntlFlock(l.file.fh.Fd(), blockLockTry, &lock); err == nil {
		l.flag ^= flag
		return true
	}
	return false
}
