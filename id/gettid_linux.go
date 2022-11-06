//go:build linux

package id

import "syscall"

func Gettid() uint64 {
	return uint64(syscall.Gettid())
}
