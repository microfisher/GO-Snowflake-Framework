package comm

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/collection"
)

// ContextDone 判断ctx是否结束
func ContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// ContextTimer 延时
func ContextTimer(ctx context.Context, delay time.Duration) (isDone bool) {
	timer := time.NewTimer(delay)
	select {
	case <-ctx.Done():
		isDone = true
	case <-timer.C:
	}
	timer.Stop()
	return

}

// NewUUIDString .
func NewUUIDString() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// FloatStringFormat 浮点字符串格式化，去除尾0和.
func FloatStringFormat(val float64) string {
	fStr := fmt.Sprintf("%f", val)
	if strings.Index(fStr, ".") >= 0 {
		for strings.HasSuffix(fStr, "0") {
			fStr = strings.TrimSuffix(fStr, "0")
		}
	}
	if strings.HasSuffix(fStr, ".") {
		fStr = strings.TrimSuffix(fStr, ".")
	}
	return fStr
}

// Float2Percent  -0.1 -> -10%
func Float2Percent(val float64) string {
	return Float2Str(val*100, 2) + "%"
}

// Float2Str .
func Float2Str(val float64, prec uint) string {
	if prec == 0 || val == 0 {
		return strconv.Itoa(int(val))
	}
	out := fmt.Sprintf("%f", val)
	dotPos := strings.Index(out, ".")
	if dotPos >= 0 {
		if len(out) > dotPos+int(prec)+1 {
			out = out[:dotPos+int(prec)+1]
		}
		out = strings.TrimRight(out, "0")
		out = strings.TrimSuffix(out, ".")
	}
	return out
}

// AmountShortString 金额转简写428.77m  56.2b
func AmountShortString(val float64) string {
	b := val / 1e9
	if b >= 1 || b <= -1 {
		return Float2Str(b, 2) + "b"
	}
	m := val / 1e6
	if m >= 1 || m <= -1 {
		return Float2Str(m, 2) + "m"
	}
	return FloatStringFormat(val)
}

// NewTimeoutCollection .
func NewTimeoutCollection(timeoutSec int64) *TimeoutCollection {
	return &TimeoutCollection{
		safemap:    collection.NewSafeMap(),
		timeoutSec: timeoutSec,
	}
}

// TimeoutCollection .
type TimeoutCollection struct {
	safemap    *collection.SafeMap
	timeoutSec int64
}

// Push 若已存在则返回false,否则push成功，返回true
func (c *TimeoutCollection) Push(key string, createTs int64) bool {
	_, has := c.safemap.Get(key)
	if has {
		return false // 不重复插入
	}
	c.safemap.Set(key, createTs)
	return true
}

// CleanupTimeout 清理超时的key
func (c *TimeoutCollection) CleanupTimeout() int {
	var delKeys []string
	now := time.Now().Unix()
	c.safemap.Range(func(key, val interface{}) bool {
		k := key.(string)
		v := val.(int64)
		if v+c.timeoutSec < now {
			delKeys = append(delKeys, k)
		}
		return true
	})
	for _, key := range delKeys {
		c.safemap.Del(key)
	}
	return len(delKeys)
}

func GetPackageName(temp interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	strs = strings.Split(strs[len(strs)-2], "/")
	return strs[len(strs)-1]
}
