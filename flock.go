package flock

import (
	"os"
	"syscall"
)

func OpenLockFile(fileName string) LockFile { return openLockFile(fileName) }

func OpenBlockLockFile(fileName string, start, end int64) LockFile {
	return openBlockLockFile(fileName, start, end)
}

type LockFile interface {
	File() *os.File
	Lock() error
	RLock() error
	TryLock() bool
	TryRLock() bool
	Unlock() error
}

const (
	/*
		The parameters of the file of the file
		see https://linux.die.net/man/2/flock
	*/
	// readLockOperation a shared lock. More than one process may hold a shared lock for a given file at a given time.
	readLockOperation = syscall.LOCK_SH
	// waiteLockOperation an exclusive lock. Only one process may hold an exclusive lock for a given file at a given time.
	waiteLockOperation = syscall.LOCK_EX
	// unLockOperation an existing lock held by this process.
	unLockOperation = syscall.LOCK_UN
	// tryReadLockOperation | tryWaiteLockOperation need to bring this parameter when you try to lock.
	tryReadLockOperation  = syscall.LOCK_NB | readLockOperation
	tryWaiteLockOperation = syscall.LOCK_NB | waiteLockOperation
)

const (
	/*
		Parameters of the content of the file
		see https://linux.die.net/man/2/fcntl
	*/

	// use in syscall.FcntlFlock() function
	// blockLockGet On input to this call, lock describes a lock we would like to place on the file.
	// If the lock could be placed, syscall.FcntlFlock() does not actually place it, but returns syscall.F_UNLCK in the l_type field of lock
	// and leaves the other fields of the structure unchanged. If one or more incompatible locks would prevent this
	// lock being placed, then syscall.FcntlFlock() returns details about one of these locks in the l_type, l_whence, l_start,
	// and l_len fields of lock and sets l_pid to be the PID of the process holding that lock.
	blockLockGet = syscall.F_GETLK
	// blockLockSet as for syscall.F_SETLK, but if a conflicting lock is held on the file,
	// then wait for that lock to be released.
	// If a signal is caught while waiting, then the call is interrupted and returns immediately
	blockLockSet = syscall.F_SETLKW
	// blockLockTry acquire a lock (when l_type is syscall.F_RDLCK or syscall.F_WRLCK) or release a lock when l_type is syscall.F_UNLCK
	// on the bytes specified by the l_whence, l_start, and l_len fields of lock.
	// If a conflicting lock is held by another process, this call returns -1 and sets errno to EACCES or EAGAIN.
	blockLockTry = syscall.F_SETLK

	// Operation instructions for blockLockGet、 blockLockSet and blockLockTry
	// bSetReadLockOperation: acquire read lock
	bSetReadLockOperation = syscall.F_RDLCK
	// bSetWaiteLockOperation: acquire only lock
	bSetWaiteLockOperation = syscall.F_WRLCK
	// bSetUnLockOperation: release lock
	bSetUnLockOperation = syscall.F_ULOCK
)
