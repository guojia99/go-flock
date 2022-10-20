package flock

import (
	"os"
)

const defaultFileFlag = 0600

type file struct {
	fileName string
	fh       *os.File
}

func (f *file) openFile() (err error) {
	if f.fh != nil {
		return
	}
	// todo 不同操作系统可能需要不同的配置，这里先用linux的
	f.fh, err = os.OpenFile(f.fileName, os.O_CREATE|os.O_RDWR, os.FileMode(defaultFileFlag))
	return
}

func (f *file) reOpenFile() (err error) {
	if st, fErr := f.fh.Stat(); fErr == nil {
		if st.Mode()&defaultFileFlag != defaultFileFlag {
			return nil
		}
		f.closeFile()
		if err = f.openFile(); err != nil {
			return err
		}
	}
	return nil
}

func (f *file) closeFile() {
	if f.fh != nil {
		_ = f.fh.Close()
		f.fh = nil
	}
}
