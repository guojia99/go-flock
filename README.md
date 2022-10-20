# go-flock File-lock

English | [简体中文](./README.ZN.md)

- Based on the `syscall.FcntlFlock` and` syscall.Flock` to achieve global file locks and local file locks.
- Only the method of the UNIX system is implemented, Window is temporarily unavailable. If you are interested, you can submit the relevant method

## Instructions

- Introduce

```
go get github.com/guojia99/go-flock
```

- Interface

```go
type LockFile interface {
	File() *os.File
	Lock() error
	RLock() error
	TryLock() bool
	TryRLock() bool
	Unlock() error
}
```

- Use global file lock

```go
package main

import (
	"github.com/guojia99/flock"
)

func main(){
    fl := flock.OpenLockFile("file.file")
    fl.Lock()
    defer fl.Unlock()
    fl.File().Write([]byte{1, 2, 3})
}
```

- Use the file local lock, of which 0, 1000 means that the lock is valid for the 0 to 1000 bytes

```go
package main

import (
	"github.com/guojia99/flock"
)

func main(){
    fl := flock.OpenBlockLockFile("file.file", 0, 1000)
    fl.Lock()
    defer fl.Unlock()
	fl.File().WriteAt([]byte{1, 2, 3}, 0)
    
    fl2 := flock.OpenBlockLockFile("file.file", 1000, 2000)
    fl2.Lock()
    defer fl2.Unlock()
    fl.File().WriteAt([]byte{1, 2, 3}, 1000)
}
```



## **License**

Go-Flock is based on the Apache 2.0 license and checks[LICENSE](./LICENSE)to get more information.
