package comm

import (
	"sync/atomic"
	"time"

	"snowflake/log"
)

// NewDelayRecorder 延时记录器
func NewDelayRecorder(prefixLog string) *DelayRecorder {
	d := &DelayRecorder{
		prefixLog: prefixLog,
	}
	d.Reset()
	return d
}

// DelayRecorder 延时记录，并发安全
type DelayRecorder struct {
	prefixLog string
	resetTS   int64
	// 单位微秒
	count, sumDelay uint64
	max, min        uint64 // 最大延时，最小延时
}

// Add 添加一个记录
func (d *DelayRecorder) Add(start time.Time) {
	d.AddDuration(time.Now().Sub(start))
}

// AddDuration 记录延时时长
func (d *DelayRecorder) AddDuration(dur time.Duration) {
	if dur < 0 {
		return
	}
	mic := uint64(dur / time.Microsecond)
	atomic.AddUint64(&d.count, +1)
	atomic.AddUint64(&d.sumDelay, mic)
	if mic > d.max {
		atomic.StoreUint64(&d.max, mic)
	}
	if mic < d.min {
		atomic.StoreUint64(&d.min, mic)
	}
}

// Reset .
func (d *DelayRecorder) Reset() {
	atomic.StoreUint64(&d.sumDelay, 0)
	atomic.StoreUint64(&d.count, 0)
	atomic.StoreUint64(&d.max, 0)
	atomic.StoreUint64(&d.min, 1<<40)
	d.resetTS = time.Now().Unix()
}

// Log .
func (d *DelayRecorder) Log(reset bool) {
	sum, count := d.sumDelay, d.count
	if count == 0 {
		return
	}
	avgMs := float64(sum) / 1000 / float64(count)
	minMs := d.min / 1000
	maxMs := d.max / 1000
	startTs := d.resetTS
	if reset {
		d.Reset()
	}
	startStr := time.Unix(startTs, 0).Format(time.RFC3339)
	log.Infof("event:delay_recorder %s start:%s count:%d delay avg_ms:%.3f max_ms:%d min_ms:%d",
		d.prefixLog, startStr, count, avgMs, maxMs, minMs)
}
