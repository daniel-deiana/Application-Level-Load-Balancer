package main
import (
	"load_balancing/data"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"io"
)

/*
	Struct used to model a backend server, storing it's
	status, the port where to reach it and other 
	metadata (protocol, etc..).
*/

/*
	In go there isn't really a enum type
	But you can emulate it using constrained int type 
	like the following 
*/

type backendState int 

/*
 define constant of type backendState with incremental values
 Each constant is one of the values of our enum
*/

const (
	StateIdle backendState = iota
	StateConnected
	StateError
	StateSleeping
)

var stateMapping = map[backendState]string {
	StateIdle 			: "idle",
	StateConnected 		: "connected",
	StateError			: "error",
	StateSleeping		: "sleep",
}

/*
	We are implementing the String() method that implements the Stringer 
	interface, this is used to print our enum returning a string
*/
func (bs backendState) String() string {
	return stateMapping[bs] 
}

type backendServer struct {
	state 		backendState
	host 		string
	protocols 	[]string 
}


func newBackendServer(newState backendState, newHost string, newProtocols []string) backendServer {
	return backendServer{state : newState, host : newHost, protocols : newProtocols}
}

/*	
	Selects the next server based on a simple round robin algorithm.
*/
func GetNextServerRR(servers map[int]backendServer, current int) (serverHost string, new_current int) {
	server := servers[current]
	new_current = (current + 1) % len(servers)
	fmt.Printf("RR: THE HOST IM RETURNING IS ", server.host)
	return server.host, new_current
}

func main () {
	/*
		I want to have a go server that has a list of other urls 
		(the actual http server to offload requests to)
		Simply start by just implmenting a round robin policy
	*/
	/*
		Istantiating the Load balancer and dispatching 
		http requests to backend servers 
	*/
	num_backends := 0 // number of backends registered
	server_index := 0
	backends := make(map[int]backendServer)
	
	s := &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	/*
		Backend registering API, The backend sends a registering message with 
		the port number where it will be listening for requests coming from the
		load balancer. 
	*/
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		buf, _ := io.ReadAll(r.Body)
		fmt.Printf("the buf contains %s\n", buf)
		var backendMsg data.BackendRegMessage		
		err := json.Unmarshal(buf,&backendMsg)
		if (err != nil) { 
			fmt.Printf("Error code of Unmarshalling is %d", err)
		}	
		fmt.Printf("the backend server sent a registration request containing %s\n", backendMsg.Host)
		serverProtocols := []string{"http"}
		backends[num_backends] = newBackendServer(StateConnected, backendMsg.Host, serverProtocols)
		fmt.Printf("Registered new backend, key is %d and host is %s", num_backends, backendMsg.Host)
		num_backends = num_backends + 1
		fmt.Fprintf(w,"ok registered !")
	})
	
	http.HandleFunc("/lb", func(w http.ResponseWriter, r *http.Request) {
	
		/* 
			Redirect the request to a backend server
			- Select one of the active backend servers from the pool. 
			- Make a new request for the backend.
			- Redirect the response of the backend to the client.
		*/

		if (len(backends) == 0) {
			fmt.Fprint(w, "There are no backends available to process request\n")		
			return
		}

		fmt.Printf("the URL of the client req is %s\n", r.URL)
		fmt.Printf("tbe HTTP header of the quest is %s\n",r.Header)
		fmt.Printf("tbe HTTP header of the quest is %s\n",r.Host)
		fmt.Printf("the HTTP scheme of the new req is %s\n\n",r.URL.Scheme)

		new_req := r.Clone(r.Context())
		new_req.RequestURI = ""
 
		var next_server string

		new_req.URL.Scheme = "http"
		next_server, new_index := GetNextServerRR(backends, server_index)
		server_index = new_index
		new_req.URL.Host = next_server
		new_req.Host = next_server

		fmt.Println("the next server is %s", next_server)
		fmt.Println("the server index is %d", server_index)

		fmt.Printf("the URL of the new req is %s\n", new_req.URL)
		fmt.Printf("tbe HTTP header of the new req is %s\n",new_req.Header)
		fmt.Printf("tbe HTTP header of the new req is %s\n",new_req.Host)
		fmt.Printf("the HTTP scheme of the new req is %s\n",new_req.URL.Scheme)
		fmt.Printf("the HTTP path of the new req is %s\n", new_req.URL.Path)

		// send the req to backend
		resp, err := http.DefaultClient.Do(new_req)
		fmt.Println("the response from the backend is %s", resp)
		fmt.Println("the error code from the backend server is %s", err)

		fmt.Fprint(w, "Hello from the backend server!")		
	})
	s.ListenAndServe()
}


