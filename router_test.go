package sleapi

import (
	"fmt"
	"net/http"
	"testing"
)

type HomeController struct {
	Controller
}

func (this *HomeController) Name() string {
	return "Home"
}

func (this *HomeController) Get(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Home Index...")
}

type AuthController struct {
	Controller
}

func (this *AuthController) Get(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Auth Index...")
}

func (this *AuthController) GetAuth(w http.ResponseWriter, req *http.Request, service string) {
	fmt.Println("Get Auth...", service)
}

func TestNewRouter(t *testing.T) {
	/*r := NewRouter()
	fmt.Println("Router: ", r.Routes)

	hc := &HomeController{}
	r.AddRoute("home", "api/home", hc)
	for k, v := range r.Routes.Table {
		fmt.Println("Key: ", k)
		fmt.Println("Value: ", v.Pattern)
		fmt.Println("Controller: ", v.Controller.Name())
		fmt.Println("Parameter: ", v.Parameter)
	}
	*/
}
