package main

import (
	"log"
	"time"

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
	return t.v
}

func (t *Test) Set(i uint64) {
	t.v = i
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

			cas_tries := 0
			slow_incr := func(i interface{}) interface{} {
				tt := i.(*Test)

				cas_tries++
				time.Sleep(100 * time.Millisecond)
				return NewTest(tt.Value() + 1)
			}

			<-do_ch // Ready, go!
			ref := tt_ref.CasUpdate(slow_incr)
			if ref == nil {
				log.Printf("Cancel(id:%d)\n", id)
				return
			}

			tt := ref.(*Test)
			log.Printf("Success(id:%d, tries:%d, value:%+v)\n",
				id, cas_tries, tt)
		}(m.New(), i)
	}

	close(do_ch)
	time.Sleep(350 * time.Millisecond)
}
