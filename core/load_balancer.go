package main

import (
	"load_balancer/datamodel"
	"load_balancer/algorithm"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"io"
)

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
	var num_backends = 0// number of backends registered
	var p *int = &num_backends

	// index of the backend to retrieve from map
	var index int = 0
	var i *int = &index  

	backends := make(map[int]datamodel.BackendServer)
	
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
		var backendMsg datamodel.BackendRegMessage		
		err := json.Unmarshal(buf,&backendMsg)

		if (err != nil) { 
			fmt.Printf("Error code of Unmarshalling is %d", err)
		}	
		
		fmt.Printf("the backend server sent a registration request containing %s\n", backendMsg.Host)
		serverProtocols := []string{"http"}
		backends[*p] = datamodel.NewBackendServer(datamodel.StateConnected, backendMsg.Host, serverProtocols)
		fmt.Printf("Registered new backend, key is %d and host is %s", num_backends, backendMsg.Host)
		*p = *p + 1
		fmt.Fprintf(w,"ok registered !")
	})

	/* 
	Redirect the request to a backend server
	- Select one of the active backend servers from the pool. 
	- Make a new request for the backend.
	- Redirect the response of the backend to the client.
	*/
	http.HandleFunc("/lb", func(w http.ResponseWriter, r *http.Request) {
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

		backendHost := algorithm.GetNextServerRR(backends, i)
		
		new_req.Header.Set("X-Forwarded-For", r.RemoteAddr)
		new_req.URL.Host = backendHost
		new_req.Host = backendHost
		new_req.URL.Scheme = "http"
		
		fmt.Println("the next server is %s", backendHost)
		fmt.Println("the server index is %d", *i)

		fmt.Printf("the URL of the new req is %s\n", new_req.URL)
		fmt.Printf("tbe HTTP header of the new req is %s\n",new_req.Host)
		fmt.Printf("the HTTP scheme of the new req is %s\n",new_req.URL.Scheme)
		fmt.Printf("the HTTP path of the new req is %s\n", new_req.URL.Path)

		// send the req to backend
		resp, err := http.DefaultClient.Do(new_req)
		fmt.Println("the response from the backend is %s", resp)
		fmt.Println("the error code from the backend server is %s", err)

		var content, cErr = io.ReadAll(resp.Body)
		if (cErr != nil) {
			panic("Load balancer: error reading backend response")
		}

		/* 
		copy the response header from the backend into the response for the client
		*/
		
		for key, value := range r.Header {
			
			// Nel video viene fatto in maniera diversa, capire perche
			w.Header().Set(key,value)
		}

		fmt.Printf("the response header from the backend said that dio is %s", w.Header().Get("dio"))

		w.Write(content)
	})
	s.ListenAndServe()
}
