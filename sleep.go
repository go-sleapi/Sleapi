package sleep

import (
	"fmt"
	//"io"
	"net/http"
	"os"
)

type Napper func(w http.ResponseWriter, req *http.Request)

type Sleep struct {
	Port   string
	Router *Router
	Naps   []Napper
}

func Sleeper() *Sleep {
	r := NewRouter()
	s := &Sleep{Port: "3030", Router: r}
	s.Naps = make([]Napper, 0)

	s.Naps = append(s.Naps, Static("/static/"))
	s.Naps = append(s.Naps, r.FindRoute)

	return s
}

// hello world, the web server
func (this *Sleep) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//io.WriteString(w, "hello, world!\n")
	fmt.Println("Server: ", req.URL.Path)

	for _, h := range this.Naps {
		h(w, req)
	}

	//this.Router.FindRoute(w, req)
}

func (this *Sleep) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3030"
	}

	http.ListenAndServe(":"+port, this)
}
