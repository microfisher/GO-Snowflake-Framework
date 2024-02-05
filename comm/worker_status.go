package comm

import (
	"context"
	"fmt"
)

// NewWorkerStatus 工作协程状态
func NewWorkerStatus(workerName string) *WorkerStatus {
	return &WorkerStatus{
		workerName: workerName,
		ExitDone:   true,
		IsReady:    false,
	}
}

type WorkerStatus struct {
	workerName string // 不变字段

	ctx      context.Context
	cancel   context.CancelFunc
	ExitDone bool   // 协程已退出完成
	IsReady  bool   // 启动加载完成,进入可用状态
	runCount uint32 // 第n次启动运行
}

func (w *WorkerStatus) Log() string {
	return fmt.Sprintf("workerName:%s runCount:%d", w.workerName, w.runCount)
}

func (w *WorkerStatus) Reset(ctx context.Context, cancel context.CancelFunc) {
	w.ctx = ctx
	w.cancel = cancel
	w.ExitDone = false // 认为开始运行了
	w.IsReady = false
	w.runCount++
}

func (w *WorkerStatus) Cancel() {
	if w.cancel != nil {
		w.cancel()
	}
}

func (w *WorkerStatus) Vaild() bool {
	return w.ctx != nil && !ContextDone(w.ctx)
}

func (w *WorkerStatus) Runing() bool {
	return w.ctx != nil && !w.ExitDone && !ContextDone(w.ctx)
}

func (w *WorkerStatus) AllReady() bool {
	if w.ExitDone || !w.IsReady || w.ctx == nil {
		return false
	}
	return !ContextDone(w.ctx)
}
