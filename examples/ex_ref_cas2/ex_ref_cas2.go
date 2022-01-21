package main

import (
	"log"
	"time"

	"github.com/l4go/task"
	"github.com/l4go/weak_ref"
)

func main() {
	m := task.NewMission()
	defer m.Done()

	var val int = 0
	w_ref := weak_ref.New(&val)

	defer log.Println("Reset")
	defer w_ref.Reset()

	do_ch := make(chan struct{})

	for i := uint64(1); i <= 5; i++ {
		go func(wm *task.Mission, id uint64) {
			defer wm.Done()

			cas_tries := 0
			<-do_ch // Ready, go!

		cas_again:
			o_ref := w_ref.Get()
			if o_ref == nil {
				log.Printf("Cancel(id:%d)\n", id)
				return
			}

			time.Sleep(100 * time.Millisecond)

			n_val := *(o_ref.(*int)) + 1
			if !w_ref.CompareAndSwap(o_ref, &n_val) {
				goto cas_again
			}

			log.Printf("Success(id:%d, tries:%d, value:%d(%+v))\n",
				id, cas_tries, n_val, &n_val)
		}(m.New(), i)
	}

	close(do_ch)
	time.Sleep(350 * time.Millisecond)
}
