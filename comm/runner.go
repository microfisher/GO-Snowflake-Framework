package comm

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"snowflake/log"
)

// Runable 可运行、停止的服务
type Runable interface {
	Start()
	Stop()
	Done() <-chan struct{}
}

// Run 运行
// 运行服务，并检测系统中断
func Run(r Runable) {
	csignal := make(chan os.Signal, 1)
	signal.Notify(csignal, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	for {
		// 运行传入的startFn
		go r.Start()
		// 等待结束信号
		var s os.Signal
		var doneQuit bool
		select {
		case s = <-csignal:
		case <-r.Done():
			s = syscall.SIGQUIT
			doneQuit = true
		}
		// close(config.StopSignal)
		switch s {
		case syscall.SIGINT:
			log.Warnf("event:service msg:Process exit with SIGINT!")
		case syscall.SIGQUIT:
			if doneQuit {
				log.Warnf("event:service Process exit with Runable Done")
			} else {
				log.Warnf("event:service process exit with SIGQUIT!")
			}
		case syscall.SIGHUP:
			time.Sleep(time.Second)
			r.Stop()
			log.Warnf("event:service process restart with SIGHUP...")
			continue // continue可以重启服务
		case syscall.SIGKILL:
			log.Warnf("event:service process killed with SIGKILL!")
		case syscall.SIGTERM:
			log.Warnf("event:service process exit with SIGTERM!")
		default:
			log.Warnf("event:service process unknown exit")
		}
		// 收到结束信号
		r.Stop()
		// 运行至此处代表关闭服务
		break
	}
}

// Runs 运行
// 运行服务，并检测系统中断
func Runs(ctx context.Context, rs ...Runable) {
	csignal := make(chan os.Signal, 1)
	signal.Notify(csignal, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	start := func() (cancelFn context.CancelFunc) {
		var wg sync.WaitGroup
		subCtx, cancel := context.WithCancel(ctx)
		for _, r := range rs {
			wg.Add(1)
			go func(r Runable) {
				defer wg.Done()
				r.Start()
				select {
				case <-subCtx.Done():
				case <-r.Done():
					log.Warnf("event:service Process exit with Runable Done")
				}
				r.Stop()
			}(r)
		}
		return func() {
			cancel()
			wg.Wait()
		}
	}
	for {
		// 运行传入的startFn
		stop := start()
		// 等待结束信号
		select {
		case s := <-csignal:
			// close(config.StopSignal)
			switch s {
			case syscall.SIGINT:
				log.Warnf("event:service Process exit with SIGINT!")
			case syscall.SIGQUIT:
				log.Warnf("event:service process exit with SIGQUIT!")
			case syscall.SIGHUP:
				stop()
				log.Warnf("event:service process restart with SIGHUP...")
				continue // continue可以重启服务
			case syscall.SIGKILL:
				log.Warnf("event:service process killed with SIGKILL!")
			case syscall.SIGTERM:
				log.Warnf("event:service process exit with SIGTERM!")
			default:
				log.Warnf("event:service process unknown exit")
			}
			// 收到结束信号
			stop()
			// 运行至此处代表关闭服务
			return
		case <-ctx.Done():
			// context要求退出
			log.Warnf("event:service context is done, stop server...")
			stop()
			return
		}
	}
}

type sigDone struct {
}

func (s sigDone) String() string {
	return "done"
}
func (s sigDone) Signal() string {
	return "done"
}

// GetCallerStackLog 调用栈日志
func GetCallerStackLog() (stacktrace string) {
	for i := 1; ; i++ {
		_, f, l, got := runtime.Caller(i)
		if !got {
			break
		}
		stacktrace += fmt.Sprintf("%s:%d\n", f, l)
	}
	return
}
