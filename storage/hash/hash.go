package hash

import (
	"fmt"
	// "github.com/metakeule/context"
	"reflect"
	"sync"
)

type (
	hash struct {
		*sync.RWMutex
		m map[uintptr]map[reflect.Type]interface{}
	}
)

func New() (ø *hash) {
	ø = &hash{&sync.RWMutex{}, map[uintptr]map[reflect.Type]interface{}{}}
	return
}

func (ø *hash) Len() int {
	return len(ø.m)
}

func (ø *hash) Inspect() string {
	return fmt.Sprintf("%v", ø.m)
}

// i should be a struct, not a pointer
func (ø *hash) Set(id uintptr, key reflect.Type, i interface{}) error {
	/*
	   if meta.Meta.Is(reflect.Ptr, i) {
	       return fmt.Errorf("pointer can't be set in context: %T", i)
	   }
	*/
	s := ø.m[id]
	if s == nil {
		s = map[reflect.Type]interface{}{}
	}
	s[key] = i
	ø.m[id] = s
	// fmt.Printf("store: %v\n", ø.m)
	return nil
}

func (ø *hash) Remove(ptr uintptr) { delete(ø.m, ptr) }

func (ø *hash) Get(id uintptr, key reflect.Type) (i interface{}, found bool) {
	// fmt.Printf("store: %v\n", ø.m)
	s := ø.m[id]
	if s == nil {
		return nil, false
	}
	// since there are no pointers allowed, i will be a copy and
	// is safe to return
	i, found = s[key]
	return
}
