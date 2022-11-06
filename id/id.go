package id

import (
	"os"
	"sync"
	"time"
)

var (
	mu                   = new(sync.Mutex)
	timestamp            = time.Now().Unix()
	deviceId             = parseIntTable[os.Getenv("DEVICE_ID")]
	processId            = os.Getpid()
	tinyThreadCounterMap = make(TinyThreadCounterMap, 0x80)
)

const offset = 1740

// Gen
// * *********** **** ***** ***** ****** ****** **** ******* ******* ********
// - ----------- ---- ----- ----- ------ ------ ---- ------- ------- --------
// - year        mon  day   hour  minute second dev  process thread  counter
func Gen() (v int64) {
	threadId := Gettid()

	mu.Lock()
	defer mu.Unlock()

	counter := tinyThreadCounterMap.Get(threadId)

PROGRESS:
	for {
		current := time.Now()

		if currentTs := current.Unix(); timestamp == currentTs {
			if *counter == 0xff {
				continue PROGRESS
			}
			*counter++
		} else {
			*counter = 1
			timestamp = currentTs
		}

		y, m, d := current.Date()
		if y-offset >= 0x800 {
			panic("id: date overflows")
		}

		v &= 0
		v <<= 11
		v |= int64(y) - offset
		v <<= 4
		v |= int64(m)
		v <<= 5
		v |= int64(d)
		v <<= 5
		v |= int64(current.Hour())
		v <<= 6
		v |= int64(current.Minute())
		v <<= 6
		v |= int64(current.Second())
		v <<= 4
		v |= int64(deviceId) & 0xf
		v <<= 7
		v |= int64(processId) & 0x7f
		v <<= 7
		v |= int64(threadId) & 0x7f
		v <<= 8
		v |= int64(*counter)
		v <<= 0
		v &= int64((1 << 63) - 1)
		break PROGRESS
	}

	return
}

type TinyThreadCounterMap map[uint8]*uint8

func (m *TinyThreadCounterMap) Get(threadId uint64) *uint8 {
	tinyThreadId := uint8(threadId & 0x7f)

	if v, ok := (*m)[tinyThreadId]; ok {
		return v
	}

	var v uint8
	(*m)[tinyThreadId] = &v

	return &v
}
