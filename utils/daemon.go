package utils

import (
	"context"
	"flag"
	"fmt"
	"github.com/sevlyar/go-daemon"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type DaemonOption struct {
	OnStop   func(cancelFunc context.CancelFunc) // cancelFunc is related to context to Run
	OnReload func()
}

func WithDaemonStop(f func(cancelFunc context.CancelFunc)) func(o *DaemonOption) {
	return func(o *DaemonOption) {
		o.OnStop = f
	}
}

func WithDaemonReload(f func()) func(o *DaemonOption) {
	return func(o *DaemonOption) {
		o.OnReload = f
	}
}

func NewDaemonOption() *DaemonOption {
	return &DaemonOption{
		OnStop: func(cancelFunc context.CancelFunc) {
			cancelFunc()
		},
	}
}

var (
	daemonMode   bool
	daemonSignal string
)

func init() {
	flag.BoolVar(&daemonMode, "d", false, "daemon mode")
	flag.StringVar(&daemonSignal, "s", "", "Send signal to the daemon: stop,reload")
}

func NormalDaemon(run func(context.Context), options ...func(o *DaemonOption)) {
	flag.Parse()

	path, _ := os.Executable()
	path = filepath.Dir(path)

	name := filepath.Base(os.Args[0])
	cntxt := &daemon.Context{
		PidFileName: filepath.Join(path, fmt.Sprintf("%s.pid", name)),
		PidFilePerm: 0644,
		LogFileName: filepath.Join(path, fmt.Sprintf("%s.out", name)),
		LogFilePerm: 0640,
		WorkDir:     path,
		Umask:       002,
	}

	ctx, cancel := context.WithCancel(context.Background())

	option := NewDaemonOption()
	for _, f := range options {
		f(option)
	}
	if option.OnStop != nil {
		daemon.AddCommand(daemon.StringFlag(&daemonSignal, "stop"), syscall.SIGTERM, func(sig os.Signal) (err error) {
			fmt.Println("stopping...")
			option.OnStop(cancel)
			return daemon.ErrStop
		})
	}

	if option.OnReload != nil {
		daemon.AddCommand(daemon.StringFlag(&daemonSignal, "reload"), syscall.SIGHUP, func(sig os.Signal) (err error) {
			fmt.Println("reloading...")
			option.OnReload()
			return nil
		})
	}

	// client process signal first
	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			fmt.Println("Unable send signal to the daemon: ", err)
			return
		}
		daemon.SendCommands(d)
		return
	}

	// start daemon
	if daemonMode {
		d, err := cntxt.Reborn()
		if err != nil {
			fmt.Println("Unable to run: ", err)
			os.Exit(-1)
		}
		if d != nil {
			return
		}
		defer cntxt.Release()

		// daemon process signal
		go func() {
			err := daemon.ServeSignals()
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}()
	} else {
		OnCloseSignal(cancel)
	}

	fmt.Println(">>> ", time.Now().Format("2006-01-02 15:04:05"))
	run(ctx)
	fmt.Println(">>> bye bye")
}

func IsDaemonMode() bool {
	return daemonMode
}
