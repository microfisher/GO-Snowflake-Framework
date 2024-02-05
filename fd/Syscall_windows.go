//go:build windows

package fd

func IncreaseFDLimit() {
	// Windows has no file descriptor limit. Do nothing.
}
