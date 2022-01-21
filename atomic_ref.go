package weak_ref

import (
	"reflect"
	"sync/atomic"
	"unsafe"
)

func create_ref(val interface{}) interface{} {
	val_v := reflect.ValueOf(val)
	if val_v.Kind() != reflect.Ptr {
		panic("unaddreable value")
	}

	addr_v := reflect.New(val_v.Type())
	addr_v.Elem().Set(val_v)
	return addr_v.Interface()
}

func to_val_ptr(typ reflect.Type, val interface{}) unsafe.Pointer {
	if val == nil {
		return nil
	}

	val_v := reflect.ValueOf(val)
	if !val_v.IsValid() {
		panic("invalid value")
	}

	if val_v.Kind() != reflect.Ptr {
		panic("unaddreable value")
	}
	if val_v.Type() != typ {
		panic("addr and val type mismatch")
	}

	return unsafe.Pointer(val_v.Pointer())
}

func to_val(typ reflect.Type, ptr unsafe.Pointer) interface{} {
	if ptr == nil {
		return nil
	}

	val_v := reflect.NewAt(typ.Elem(), ptr)
	if !val_v.IsValid() {
		return nil
	}

	return val_v.Interface()
}

func to_addr_ptr(addr interface{}) (reflect.Type, unsafe.Pointer) {
	addr_v := reflect.ValueOf(addr)

	if !addr_v.IsValid() {
		panic("invalid value")
	}
	if addr_v.Kind() != reflect.Ptr {
		panic("unaddreable value")
	}
	if addr_v.Type().Elem().Kind() != reflect.Ptr {
		panic("unaddreable value")
	}

	return addr_v.Type(), unsafe.Pointer(addr_v.Pointer())
}

func cas_ref(addr interface{}, old_val, new_val interface{}) bool {
	addr_typ, addr_ptr := to_addr_ptr(addr)
	swap := atomic.CompareAndSwapPointer((*unsafe.Pointer)(addr_ptr),
		to_val_ptr(addr_typ.Elem(), old_val),
		to_val_ptr(addr_typ.Elem(), new_val))
	return swap
}

func swap_ref(addr interface{}, val interface{}) interface{} {
	addr_typ, addr_ptr := to_addr_ptr(addr)
	old_ptr := atomic.SwapPointer((*unsafe.Pointer)(addr_ptr),
		to_val_ptr(addr_typ.Elem(), val))
	return to_val(addr_typ.Elem(), old_ptr)
}

func store_ref(addr interface{}, val interface{}) {
	addr_typ, addr_ptr := to_addr_ptr(addr)
	atomic.StorePointer((*unsafe.Pointer)(addr_ptr),
		to_val_ptr(addr_typ.Elem(), val))
}

func load_ref(addr interface{}) interface{} {
	addr_typ, addr_ptr := to_addr_ptr(addr)
	ptr := atomic.LoadPointer((*unsafe.Pointer)(addr_ptr))
	return to_val(addr_typ.Elem(), ptr)
}
