package initialize

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Option for logrus
const logPath = "logs"

type Option struct {
	Level     logrus.Level
	Pid       bool
	ToConsole bool
	Name      string
}

func NewOption() *Option {
	return &Option{
		Level:     logrus.InfoLevel,
		Pid:       false,
		ToConsole: false,
		Name:      filepath.Base(os.Args[0]),
	}
}
func WithLevel(l string) func(o *Option) {
	level, err := logrus.ParseLevel(l)
	if err != nil {
		panic(err)
	}

	return func(o *Option) {
		o.Level = level
	}
}

func WithToConsole() func(o *Option) {
	return func(o *Option) {
		o.ToConsole = true
	}
}

func WithDaemon(daemon bool) func(o *Option) {
	return func(o *Option) {
		o.ToConsole = !daemon
	}
}

func WithPid() func(o *Option) {
	return func(o *Option) {
		o.Pid = true
	}
}

func WithName(name string) func(o *Option) {
	return func(o *Option) {
		o.Name = name
	}
}

func getLogFile(name string, withpid bool) io.Writer {
	p, err := os.Executable()
	if err != nil {
		panic(err)
	}
	root := filepath.Dir(p)

	_ = os.Mkdir(filepath.Join(root, logPath), os.ModePerm)

	name = fmt.Sprintf("%s.log", name)
	if withpid {
		name = fmt.Sprintf("%s_%d.log", name, os.Getpid())
	}

	logName := filepath.Join(root, logPath, name)

	// https://github.com/lestrrat-go/strftime
	logf, err := rotatelogs.New(
		logName+".%y%m%d",
		rotatelogs.WithLinkName(logName),
		rotatelogs.WithMaxAge(30*24*time.Hour),    // 7 days
		rotatelogs.WithRotationTime(24*time.Hour), //  1 day
	)
	if err != nil {
		panic(err)
	}

	return logf
}

// HookFilename -
type HookFilename struct {
}

// Fire -
func (h *HookFilename) Fire(entry *logrus.Entry) error {
	if entry.HasCaller() {
		entry.Caller.File = filepath.Base(entry.Caller.File)
	}
	return nil
}

// Levels -
func (h *HookFilename) Levels() []logrus.Level {
	return logrus.AllLevels
}

func LoggerInit(options ...func(o *Option)) {
	op := NewOption()
	for _, o := range options {
		o(op)
	}

	logrus.SetReportCaller(true)
	logrus.SetLevel(op.Level)

	if !op.ToConsole {
		writer := getLogFile(op.Name, op.Pid)
		logrus.SetOutput(writer)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006/01/02 15:04:05.000",
	})

	logrus.AddHook(&HookFilename{})
}
