package comm

import (
	"strconv"
	"strings"
	"unsafe"
)

// Str2bytes 避免内存复制
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	b := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}

// Bytes2str 避免内存复制
func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// UnescapeUnicode unicode字符串解析
// "hello \\u4f60\\u597d" --> "hello 你好"
func UnescapeUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
