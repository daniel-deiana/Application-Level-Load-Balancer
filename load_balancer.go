package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"io"
)

/*	
	Selects the next server based on a simple round robin algorithm.
*/

type backendRegMessage struct{
	Port string `json:"Port"`
}

func GetNextServerRR(servers map[int]string, current int) (selected_server string, new_current int) {
	selected_server = servers[current]
	current = (current + 1) % len(servers)
	return selected_server, new_current
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

	backends := make(map[int]string)
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
		var backendMsg backendRegMessage
		
		err := json.Unmarshal(buf,&backendMsg)

		if (err != nil) { 
			fmt.Printf("Error code of Unmarshalling is %d", err)
		}
		fmt.Printf("the backend server sent a registration request containing %s\n", backendMsg.Port)
		host := "localhost:" + backendMsg.Port
		backends[num_backends] = host
		fmt.Printf("Registered new backend, key is %d and host is %s", num_backends, host)
		num_backends = num_backends + 1
		fmt.Fprintf(w,"ok registered !")
	})
	
	http.HandleFunc("/lb", func(w http.ResponseWriter, r *http.Request) {
	
		/* 
			Redirect the request to a backend server
			- Pool one of the active backend servers from the pool. 
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


