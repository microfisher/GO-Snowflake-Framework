package comm

import (
	"fmt"
	"runtime"
)

// 获取系统信息
func GetSystemInfo() (memory uint64, goroutine int) {
	var systemMemory runtime.MemStats
	runtime.ReadMemStats(&systemMemory)
	memory = systemMemory.Sys / 1024 / 1024
	goroutine = runtime.NumGoroutine()
	return
}

// GetProcessStatusLog 进程状态打印信息
func GetProcessStatusLog() string { // 影响进程全局的耗时操作,调用间隔至少一分钟以上
	var systemMemory runtime.MemStats
	runtime.ReadMemStats(&systemMemory)
	heapMem := systemMemory.Alloc / 1024 / 1024
	stackMem := systemMemory.StackSys / 1024 / 1024
	sysMem := systemMemory.Sys / 1024 / 1024
	goroutine := runtime.NumGoroutine()
	return fmt.Sprintf("event:process_status heap_mem_mb:%d stack_mem_mb:%d sys_mem_mb:%d goroutine:%d",
		heapMem, stackMem, sysMem, goroutine)
}
