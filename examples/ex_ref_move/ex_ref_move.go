package main

import (
	"log"
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
	defer tt_ref.Reset()

	go free_worker(m.New(), tt_ref.Move())

	if tt_ref.Get() == nil {
		log.Println("No cleaning required.")
		return
	}
	log.Println("Cleaning.")
}

func free_worker(wm *task.Mission, tt_ref *weak_ref.WeakRef) {
	defer wm.Done()
	defer tt_ref.Reset()

	ref := tt_ref.Get()
	if ref == nil {
		log.Println("Fail worker")
		return
	}
	log.Println("Success worker")
}
