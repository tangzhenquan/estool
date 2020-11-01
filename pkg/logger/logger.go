package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type ConfigT struct {
	File  string `validate:"required"`
	Level string `validate:"required"`
}

func InitLogger(c *ConfigT) error {
	if ok := strings.HasPrefix(c.File, "/"); !ok {
		dir := filepath.Dir(c.File)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
	}

	writer, err := rotatelogs.New(
		c.File+".%Y%m%d",
		rotatelogs.WithLinkName(c.File),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(7),
	)
	if err != nil {
		return err
	}

	level, err := logrus.ParseLevel(c.Level)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	var writeMap = make(lfshook.WriterMap)
	for _, l := range logrus.AllLevels {
		if l <= level {
			writeMap[l] = writer
		}
	}
	logrus.AddHook(lfshook.NewHook(
		writeMap,
		&logrus.JSONFormatter{},
	))
	return nil
}
