package utils

import (
	"github.com/sirupsen/logrus"
	"runtime"
)

func SafeExecFunc(fn func(...interface{}), args ...interface{}) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithFields(logrus.Fields{
					"err": err,
				}).Error("Bug!!!, safe exec function error, please check your code, stack=", string(Stacks()))
			}
		}()

		fn(args...)
	}()
}

func Stacks() []byte {
	const size = 16 << 10
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	return buf
}
