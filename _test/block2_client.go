/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/2/26 下午5:22.
 *  * Author: guojia(https://github.com/guojia99)
 */

package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/guojia99/flock"
)

func main() {
	if len(os.Args) < 3 {
		log.Println("you will run `go run client.go client1 100`")
		return
	}
	file := "client.test"
	clientName := os.Args[1]
	start, _ := strconv.Atoi(os.Args[2])

	fl := flock.OpenBlockLockFile(file, int64(start), int64(start+100))
	test := func() {
		log.Printf("%s try lock -> %v\n", clientName, fl.TryLock())
		_ = fl.Lock()
		log.Printf("%s lock file\n", clientName)
		defer fl.Unlock()
		time.Sleep(time.Second * 10)
	}
	test()
}
