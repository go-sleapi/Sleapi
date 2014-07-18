Sleapi
======

Lightweight Go Web Framework for building REST applications.  Modeled after ASP.NET Web API (mostly for how routes are defined).

Example:

		type MainController struct {
			sleapi.Controller
		}
		
		//This will get called for any GET requests that match the route 'api/main'
		func (this *MainController) Get(w http.ResponseWriter, req *http.Request){
		  fmt.Fprintf(w, "Hello World")
		}
		
		func (this *MainController) GetName(w http.ResponseWriter, req *http.Request, name string){
		  fmt.Fprintf(w, "Hello World, " + name)
		}
		
		func (this *MainController) PostName(w http.ResponseWriter, req *http.Request){
		  body, err := ioutil.ReadAll(req.Body)
		  defer req.Body.Close()
		  if err != nil {
		      log.Fatal(err)
		  }
		
		  //The body would typiclaly be JSON send from the client
		  fmt.Println("Body: ", string(body))
		}
		
		func main() {
			s := sleapi.Sleeper()
		
			mainController := &MainController{}
			s.Router.AddRoute("Main", "api/main", mainController)
			s.Router.AddRoute("Name", "api/{name}, mainController)
		
			s.Run()
		}
