package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sleep/sleep"
	"sleep/sync"
)

type SyncController struct {
	sleep.Controller
}

type Data struct {
	Foo string
	Bar string
	Baz int
}

func (this *SyncController) Get(w http.ResponseWriter, req *http.Request) {
	data := &Data{"Test", "Testing", 100}
	b, _ := json.Marshal(data)
	io.WriteString(w, string(b))
}

func (this *SyncController) Post(w http.ResponseWriter, req *http.Request, when string) {
	fmt.Println("When: " + when)
	/*data := &Data{"Test", "Testing", 100}
	b, _ := json.Marshal(data)
	io.WriteString(w, string(b))
	*/

	go func( /*done chan bool*/) {
		sync.EnqueueJob()
		sync.ProcessJobQueue()

		//done <- true
	}( /*done*/)
}

func main() {
	s := sleep.Sleeper()
	fmt.Println("Sleeper: ", s)
	fmt.Println("Testing...")

	s.Router.AddRoute("Sync", "app/sync/{when}", &SyncController{})

	s.Run()
}
