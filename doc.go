// Copyright 2013 Marc René Arns. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package context offers a reliabel way to share context in multiple goroutines.
You may use it to add context or sessions to http requests for example.

For the actual storage you need a storage that implements the Storage interface:

    Storage interface {
        Set(uintptr, int, interface{}) error
        Get(uintptr, Object) (found bool)
        Remove(uintptr)
        Lock()
        Unlock()
        RLock()
        RUnlock()
    }

A simple example for such a storage is provided at github.com/metakeule/context/storage/hash

Usage

    package main

    import (
        "github.com/metakeule/context"
        "github.com/metakeule/context/storage/hash"
    )

    // you typically want to save structs, keep in mind that 
    // the objects are saved with the type as key, so for
    // each area you might only have one object per type
    type User struct {FirstName, LastName  string}

    // the struct must fullfill the Object interface:
    //  - Restore replaces the current instance with the given interface
    func (ø *User) Restore(i interface{}) {
        u := i.(User)
        *ø = u
    }

    // each place where the same objects are shared has to be identified
    // by a unique pointer. here we take a string and will submit the pointer to it
    var area = "shared1"

    // create a new context. we pass a storag (hash.New()) and the number of getter
    // and setter channels
    var ctx = context.New(hash.New(), 50, 100)

    func main() {
        // start the contextmanager
        go ctx.Run()

        u := User{"Donald","Duck"}
        // save the value of user.
        // note that you have to pass the pointer of the area and
        // the object you want to save even if the object is copied on the storage
        err := ctx.Set(&area, &u)

        // somewhere else

        user := User{}
        // here again: two pointers
        found := ctx.Get(&area, &user)
        // do whatever you want with user

        // sometimes it is important to cleanup the area:
        ctx.Remove(&area)
    }

*/
package context
