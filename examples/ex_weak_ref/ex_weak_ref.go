package main

import (
	"log"
	"runtime"
	"sync/atomic"

	"github.com/l4go/task"
	"github.com/l4go/weak_ref"
)

type Test struct {
	v uint64
}

func NewTest(i uint64) *Test {
	return &Test{v: i}
}
func (t *Test) Value() uint64 {
	return atomic.LoadUint64(&t.v)
}

func (t *Test) Set(i uint64) {
	atomic.StoreUint64(&t.v, i)
}

func main() {
	m := task.NewMission()
	defer m.Done()

	tt := NewTest(0)
	tt_ref := weak_ref.New(tt)
	defer log.Println("Reset")
	defer tt_ref.Reset()

	do_ch := make(chan struct{})

	for i := uint64(1); i <= 5; i++ {
		go func(wm *task.Mission, id uint64) {
			defer wm.Done()

			<-do_ch

			ref := tt_ref.Get()
			if ref == nil {
				log.Println("Cancel:", id)
				return
			}
			tt := ref.(*Test)
			log.Println("Success:", id, tt.Value())
		}(m.New(), i)

	}

	tt.Set(1)
	close(do_ch)
	tt.Set(2)
	runtime.Gosched()
	tt.Set(3)
}
