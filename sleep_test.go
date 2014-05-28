package sleapi

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNewSleep(t *testing.T) {
	s := Sleeper()
	fmt.Println("Sleeper: ", s)
	fmt.Println("Router: ", s.Router)

	hc := &HomeController{}
	s.Router.AddRoute("home", "api/home", hc)
	s.Router.AddRoute("authorization", "api/auth/{service}", &AuthController{})
	go s.Run()

	//http.Get("http://localhost:" + s.Port + "/api/auth/strava")
	//http.Get("http://localhost:" + s.Port + "/api/home")
	//http.Get("http://localhost:" + s.Port + "/api/home/test")
	http.Get("http://localhost:" + s.Port + "/static/index.html")
}
