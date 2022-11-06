//go:build darwin && cgo

package id

// #include <stdlib.h>
// #include <pthread.h>
// uint64_t gettid() {
// 	   uint64_t tid;
//     pthread_threadid_np(NULL, &tid);
//     return tid;
// }
import "C"

func Gettid() uint64 {
	return uint64(C.gettid())
}
