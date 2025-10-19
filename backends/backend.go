package main
import (
	"fmt"
	"net/http"
	"time"
	"os"
	"bytes"
)

/*
	TODO:
	- Add a struct to model backend servers
	- Add methods to interact with backend servers, health cheking
*/

func parsePort() (port string) {
	parameterArgs := os.Args[1:]
	if (len(parameterArgs) < 1) {
		return ""
	}
	port = parameterArgs[0]
	return port 
}

func main () {

	/*
		I want to have a go server that has a list of other urls (the actual http server to offload requests to)
		Simply start by just implmenting a round robin policy
	*/

	/*
		Istantiating the Load balancer and dispatching http requests to backend servers 
	*/

	port := parsePort()

	if (port == ""){
		// user needs to specify a port
		panic(true)
	}

	fmt.Println("The port im sending to the server is %s", port)

	// Send host made as localhost:port to load_balancer
	jsonStr := "{ \"Host\" : \"" + "localhost:" + port +"\" }"
	
	// creates a Buffer type (a slice of byte) initializing it with the contents of a string
	buf := bytes.NewBufferString(jsonStr)
	
	fmt.Println("printing the byte slice created from the jsonStr %s", buf)
	_, err := http.Post("http://localhost:8080/register", "application/json", buf)

	if (err != nil) {
		fmt.Printf("error on making the register to the load balancer")
		return
	}
	
	// send register request to the load balancer 
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	http.HandleFunc("/lb", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("The X-forwarded-for is %s", r.Header.Get("X-Forwarded-For"))    	
		fmt.Println("Backend server at port %s responding to load balancer request!", port)
		var xfwdHost = r.Header.Get("X-Forwarded-Host")
		fmt.Printf("the x fowarded host of the client is %s", xfwdHost)
		fmt.Fprint(w, "Backend has responded! Hello client")		
	})
	
    http.HandleFunc("/health", func(w http.ResponseWriter, r * http.Request) {
		fmt.Fprintf(w, "Received health checking messages")
    })

	s.ListenAndServe()
}
