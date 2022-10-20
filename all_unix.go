//go:build !windows

package flock

import (
	"os"
	"sync"
	"syscall"
)

func openLockFile(fileName string) *lockFile {
	return &lockFile{file: &file{fileName: fileName}}
}

type lockFile struct {
	m    sync.Mutex
	file *file
	flag int
}

func (l *lockFile) File() *os.File { return l.file.fh }
func (l *lockFile) Lock() error    { return l.lock(waiteLockOperation | readLockOperation) }
func (l *lockFile) RLock() error   { return l.lock(readLockOperation) }
func (l *lockFile) TryLock() bool  { return l.tryLock(tryWaiteLockOperation | readLockOperation) }
func (l *lockFile) TryRLock() bool { return l.tryLock(tryReadLockOperation) }
func (l *lockFile) Unlock() error  { return l.unlock() }

func (l *lockFile) lock(flag int) (err error) {
	l.m.Lock()
	defer l.m.Unlock()

	if l.flag&flag == flag {
		return nil
	}

	if err = l.file.openFile(); err != nil {
		return err
	}

	err = syscall.Flock(int(l.file.fh.Fd()), flag)
	return
}

func (l *lockFile) unlock() error {
	l.m.Lock()
	defer l.m.Unlock()

	if l.flag == 0 || l.file.fh == nil {
		return nil
	}

	if err := syscall.Flock(int(l.file.fh.Fd()), unLockOperation); err != nil {
		return err
	}
	l.file.closeFile()
	l.flag = 0
	return nil
}

func (l *lockFile) tryLock(flag int) bool {
	l.m.Lock()
	defer l.m.Unlock()

	if l.flag&flag == flag {
		return true
	}

	if err := l.file.openFile(); err != nil {
		return false
	}

	err := syscall.Flock(int(l.file.fh.Fd()), flag)
	return err == nil
}
