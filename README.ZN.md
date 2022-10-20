# go-flock 文件锁

[English](./README.md) | 简体中文

- 基于 `syscall.FcntlFlock` 和 `syscall.Flock`实现全局文件锁和局部文件锁。
- 仅实现了unix系统的方法， window暂时无法使用，感兴趣可以提交相关方法

## 使用方法

- 引入

```
go get github.com/guojia99/go-flock
```

- 接口

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

- 使用全局文件锁

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

- 使用文件局部锁,  其中0, 1000 代表该锁仅对该0～1000字节的有效

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

## **许可证**

go-flock 基于 Apache 2.0 许可证，查看 [LICENSE](./LICENSE) 获取更多信息。









