//go:build !windows

package fd

import (
	"syscall"

	"snowflake/log"
)

// IncreaseFDLimit 增大文件描述符上限
func IncreaseFDLimit() {
	var rlm syscall.Rlimit
	// limit.Cur is the soft limit. limit.Max is the hard limit.

	// Try to increase the soft limit
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlm)
	if rlm.Cur < 65535 && rlm.Cur < rlm.Max {
		rlm.Cur = rlm.Max
		syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlm)
	}

	// Try to increase the hard limit
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlm)
	if rlm.Cur < 65535 || rlm.Max < 65535 {
		rlm.Cur = 65535
		rlm.Max = 65535
		syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlm)
	}

	// checking
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlm)
	log.Infof("File descriptor soft limit:%d hard limit:%d ", rlm.Cur, rlm.Max)
	if rlm.Max < 20000 {
		log.Error("File descriptor ard limit is too small: ", rlm.Max, "! The problem may be solved by executing the following command before launching: ulimit -Hn 65535")
	}
	if rlm.Cur < 20000 {
		log.Error("File descriptor soft limit is too small: ", rlm.Cur, "! The problem may be solved by executing the following command before launching: ulimit -Sn 65535")
	}
}
