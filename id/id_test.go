package id

import (
	"sync"
	"testing"
	"time"
)

func TestID(t *testing.T) {
	const N = 3000

	var (
		wg = new(sync.WaitGroup)
		mu = new(sync.Mutex)
		mp = make(map[int64]struct{}, 1000)
	)

	for i := 0; i < N; i++ {
		wg.Add(1)
		go test(t, wg, mu, mp)
	}

	wg.Wait()
	t.Log("finish")
}

func test(t *testing.T, wg *sync.WaitGroup, mu *sync.Mutex, mp map[int64]struct{}) {
	const format = "2006-01-02 15:04:05"

	mu.Lock()
	defer mu.Unlock()
	defer wg.Done()

	did, pid, tid := deviceId, processId, Gettid()

	var id = Gen()
	if _, exist := mp[id]; exist {
		t.Errorf("duplicate error: %d", id)
		return
	}
	mp[id] = struct{}{}

	eTime, eDeviceId, eProcessId, eThreadId, eCounter := inspect(id)
	if delta := time.Now().Sub(eTime); !(time.Second*0 <= delta && delta <= time.Second*3) {
		t.Errorf("time error: %s\n", eTime.Format(format))
		return
	}
	if eDeviceId != did&0xf {
		t.Errorf("device_id error: %d, %d\n", eDeviceId, did)
		return
	}
	if int(eProcessId) != pid&0x7f {
		t.Errorf("process_id error: %d, %d\n", eProcessId, pid&0x7f)
		return
	}
	if uint64(eThreadId) != tid&0x7f {
		t.Errorf("thread_id error: %d, %d\n", eThreadId, tid&0x7f)
		return
	}
	if eCounter != *tinyThreadCounterMap.Get(tid) {
		t.Errorf("counter error: %d, %d\n", eCounter, *tinyThreadCounterMap.Get(tid))
		return
	}

	t.Logf(
		"%d: time=%s, device_id=%d, process_id=%d, thread_id=%d, counter=%d",
		id,
		eTime.Format(format),
		eDeviceId,
		eProcessId,
		eThreadId,
		eCounter,
	)
	return
}

func inspect(id int64) (t time.Time, did uint8, pid uint8, tid uint8, c uint8) {
	c = uint8(id & 0xff)
	id >>= 8
	tid = uint8(id & 0x7f)
	id >>= 7
	pid = uint8(id & 0x7f)
	id >>= 7
	did = uint8(id & 0xf)
	id >>= 4
	S := int(id & 0x3f)
	id >>= 6
	M := int(id & 0x3f)
	id >>= 6
	H := int(id & 0x1f)
	id >>= 5
	d := int(id & 0x1f)
	id >>= 5
	m := int(id & 0xf)
	id >>= 4
	y := int(id) + offset
	t = time.Date(y, time.Month(m), d, H, M, S, 0, time.Local)
	return
}
