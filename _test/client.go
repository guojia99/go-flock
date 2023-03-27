/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/2/26 下午5:22.
 *  * Author: guojia(https://github.com/guojia99)
 */

package main

import (
	"log"
	"os"
	"time"

	"github.com/guojia99/flock"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("you will run `go run client.go client1`")
		return
	}

	clientName := os.Args[1]
	doubleFile := "double_client.test"
	fl := flock.OpenLockFile(doubleFile)

	test1 := func() {
		log.Printf("%s try lock -> %v\n", clientName, fl.TryLock())
		_ = fl.Lock()

		log.Printf("%s lock file\n", clientName)
		defer fl.Unlock()
		fl.File().WriteAt([]byte{1, 2, 3}, 0)
		time.Sleep(time.Second * 10)
	}

	test2 := func() {
		log.Printf("%s try read lock -> %v\n", clientName, fl.TryRLock())
		_ = fl.RLock()
		log.Printf("%s read lock file\n", clientName)
		defer fl.Unlock()
		time.Sleep(time.Second * 10)
	}

	test1()
	test2()
}
