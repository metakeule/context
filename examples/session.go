package main

import (
	"fmt"
	"github.com/metakeule/context"
	"github.com/metakeule/context/storage/hash"
	"net/http"
)

type sessionId struct{ name string }

var (
	store   = hash.New()
	Session = context.New(store, 50, 100)
	links   = `<br>
    <a href="/set?firstname=Bugs&lastname=Bunny&id=1&sessionid=a">Change to Bugs Bunny </a><br>
    <a href="/set?firstname=Donald&lastname=Duck&id=2&sessionid=a">Change to Donald Duck </a>
    </body></html>`
	sessions = map[string]*sessionId{}
)

func getSessionId(r *http.Request) (sess *sessionId) {
	s := r.URL.Query().Get("sessionid")
	var ok bool
	sess, ok = sessions[s]
	if !ok {
		//	fmt.Println("session not found: ", s)
		sess = &sessionId{s}
		sessions[s] = sess
	}
	return
}

type User struct{ Id, FirstName, LastName string }

func (ø *User) Restore(i interface{}) { u := i.(User); *ø = u }

func SetUser(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	user := &User{q.Get("id"), q.Get("firstname"), q.Get("lastname")}
	err := Session.Set(getSessionId(r), user)
	fmt.Printf("size of store: %v\n", store.Len())
	if err != nil {
		w.Write([]byte(`Error: ` + err.Error()))
		return
	}
	w.Write([]byte(`<html><body>Go to <a href="/show?sessionid=a">Show User</a></body></html>`))
}

func ShowUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(store.Inspect())
	user := &User{}
	found := Session.Get(getSessionId(r), user)
	if found {
		w.Write([]byte(`<html><body>Current user: ` +
			user.FirstName + " " + user.LastName + "(" + user.Id + ")" + links))
		return
	}
	w.Write([]byte(`<html><body>Current user: not found` + links))
}

func main() {
	go Session.Run()
	http.HandleFunc("/show", ShowUser)
	http.HandleFunc("/set", SetUser)
	http.ListenAndServe(":8080", nil)
}
