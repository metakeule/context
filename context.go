package context

import (
	"fmt"
	"reflect"
)

type (
	Object interface {
		Restore(interface{})
	}

	Storage interface {
		Set(uintptr, reflect.Type, interface{}) error
		Get(uintptr, reflect.Type) (obj interface{}, found bool)
		Remove(uintptr)
		Lock()
		Unlock()
		RLock()
		RUnlock()
	}

	getQuery struct {
		Ptr    uintptr
		Object Object
		Found  chan bool // found
	}

	getChan chan *getQuery

	setQuery struct {
		Ptr    uintptr
		Object Object
		Error  chan error
	}
	setChan chan *setQuery

	Context struct {
		storage     Storage
		setChannels setChan
		getChannels getChan
		Debug       bool
	}
)

func New(store Storage, numSetChannels int, numGetChannels int) (ø *Context) {
	ø = &Context{
		storage: store,
	}
	ø.setChannels = make(setChan, numSetChannels)
	ø.getChannels = make(getChan, numGetChannels)
	return
}

func (ø *Context) pointer(i interface{}) uintptr {
	return reflect.ValueOf(i).Pointer()
}

// we have an own method to be able to use defer
func (ø *Context) set(ptr uintptr, key reflect.Type, val interface{}) error {
	if ø.Debug {
		fmt.Println("setting key %v in Context %v to %#v", key, ptr, val)
	}
	ø.storage.Lock()
	defer ø.storage.Unlock()
	return ø.storage.Set(ptr, key, val)
}

// we have an own method to be able to use defer
func (ø *Context) get(ptr uintptr, obj Object) (found bool) {
	if ø.Debug {
		fmt.Println("get key %v in Context %v to %#v", reflect.TypeOf(obj), ptr)
	}
	ø.storage.RLock()
	// fmt.Println("aquired storage lock for get")
	defer ø.storage.RUnlock()
	var i interface{}
	i, found = ø.storage.Get(ptr, reflect.TypeOf(obj))
	if found {
		obj.Restore(i)
	}
	return
}

func (ø *Context) Run() {
	for {
		select {
		case g := <-ø.getChannels:
			g.Found <- ø.get(g.Ptr, g.Object)
		case s := <-ø.setChannels:
			// pass the value the object points to, to save a copy
			// an return a copy to the set method of Object, that is
			// called by getFromContext
			i := reflect.ValueOf(s.Object).Elem().Interface()
			s.Error <- ø.set(s.Ptr, reflect.TypeOf(s.Object), i)
		}
	}
}

func (ø *Context) Set(r interface{}, obj Object) (err error) {
	resp := make(chan error, 1)
	ø.setChannels <- &setQuery{
		Ptr:    ø.pointer(r),
		Object: obj,
		Error:  resp,
	}
	err = <-resp
	return
}

// val should be a pointer to a struct
func (ø *Context) Get(r interface{}, obj Object) (found bool) {
	resp := make(chan bool, 1)
	ø.getChannels <- &getQuery{
		Ptr:    ø.pointer(r),
		Object: obj,
		Found:  resp,
	}
	found = <-resp
	return
}

func (ø *Context) Remove(i interface{}) {
	ptr := ø.pointer(i)
	if ø.Debug {
		fmt.Println("removing Context %v", ptr)
	}
	ø.storage.Lock()
	defer ø.storage.Unlock()
	ø.storage.Remove(ptr)
}
