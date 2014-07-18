Sleapi
======

Lightweight Go Web Framework for building REST based ASP.NET MVC applications.

Example:

		type MainController struct {
			sleapi.Controller
		}
		
		func (this *MainController) Get(w http.ResponseWriter, req *http.Request){
		  fmt.Fprintf(w, "Hello World")
		}
		
		func main() {
			s := sleapi.Sleeper()
		
			mainController := &MainController{}
			s.Router.AddRoute("Main", "/", mainController)
		
			s.Run()
		}
