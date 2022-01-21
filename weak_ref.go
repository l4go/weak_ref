package weak_ref

type WeakRef struct {
	ptr interface{}
}

func New(val interface{}) *WeakRef {
	return &WeakRef{ptr: create_ref(val)}
}

func (wr *WeakRef) Move() *WeakRef {
	val := swap_ref(wr.ptr, nil)
	return New(val)
}

func (wr *WeakRef) Set(val interface{}) {
	store_ref(wr.ptr, val)
}

func (wr *WeakRef) Swap(new_val interface{}) interface{} {
	return swap_ref(wr.ptr, new_val)
}

func (wr *WeakRef) CompareAndSwap(old_val, new_val interface{}) bool {
	return cas_ref(wr.ptr, old_val, new_val)
}

type CasUpdateFunc func(v interface{}) interface{}

func (wr *WeakRef) CasUpdate(f CasUpdateFunc) interface{} {

cas_again:
	old_val := load_ref(wr.ptr)
	if old_val == nil {
		return nil
	}

	new_val := f(old_val)
	if !cas_ref(wr.ptr, old_val, new_val) {
		goto cas_again
	}

	return new_val
}

func (wr *WeakRef) Get() interface{} {
	now_val := load_ref(wr.ptr)
	if now_val == nil {
		return nil
	}
	return now_val
}

func (wr *WeakRef) Reset() {
	store_ref(wr.ptr, nil)
}
