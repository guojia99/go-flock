/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/2/26 下午5:22.
 *  * Author: guojia(https://github.com/guojia99)
 */

package flock

import (
	"fmt"
	"os"
	"testing"
	"time"
)

const fileNameFormat = "./test/test_file_%d.test"

func init() {
	_ = os.MkdirAll("test", os.ModeDir|os.ModePerm)
}

func Test_lockFile(t *testing.T) {
	t.Run("lock and unlock", func(t *testing.T) {
		fileName := fmt.Sprintf(fileNameFormat, time.Now().Unix())
		fl := openLockFile(fileName)
		if err := fl.Lock(); err != nil {
			t.Fatalf("lock error %s", err)
		}
		if err := fl.Unlock(); err != nil {
			t.Fatalf("unlock error %s", err)
		}
	})

	t.Run("try lock", func(t *testing.T) {
		fileName := fmt.Sprintf(fileNameFormat, time.Now().Unix())
		fl := openLockFile(fileName)
		if ok := fl.TryLock(); !ok {
			t.Fatalf("can not lock")
		}
		if err := fl.Unlock(); err != nil {
			t.Fatalf("unlock error %s", err)
		}
	})

	t.Run("double lock", func(t *testing.T) {
		fileName := fmt.Sprintf(fileNameFormat, time.Now().Unix())
		fl := openLockFile(fileName)
		defer fl.Unlock()

		if err := fl.Lock(); err != nil {
			t.Fatalf("lock error %s", err)
		}

		fl2 := openLockFile(fileName)
		if err := fl2.Lock(); err != nil {
			t.Logf("lock error %s", err)
		}
	})
}
