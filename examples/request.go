package main

import (
	"fmt"
	"github.com/metakeule/context"
	"github.com/metakeule/context/storage/hash"
	"net/http"
)

var (
	store   = hash.New()
	Context = context.New(store, 50, 100)
	links   = `<br>
    <a href="/set?firstname=Bugs&lastname=Bunny&id=1">Change to Bugs Bunny </a><br>
    <a href="/set?firstname=Donald&lastname=Duck&id=2">Change to Donald Duck </a>
    </body></html>`
)

type User struct{ Id, FirstName, LastName string }

func (ø *User) Restore(i interface{}) { u := i.(User); *ø = u }

func SetUser(w http.ResponseWriter, r *http.Request) {
	defer Context.Remove(r)
	q := r.URL.Query()
	user := &User{q.Get("id"), q.Get("firstname"), q.Get("lastname")}
	err := Context.Set(r, user)
	fmt.Printf("size of store: %v\n", store.Len())
	if err != nil {
		w.Write([]byte(`Error: ` + err.Error()))
		return
	}
	ShowUser(w, r)
}

func ShowUser(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	found := Context.Get(r, user)
	if found {
		w.Write([]byte(`<html><body>Current user: ` +
			user.FirstName + " " + user.LastName + "(" + user.Id + ")" + links))
		return
	}
	w.Write([]byte(`<html><body>Current user: not found` + links))
}

func main() {
	go Context.Run()
	http.HandleFunc("/set", SetUser)
	http.ListenAndServe(":8080", nil)
}
