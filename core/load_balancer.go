package main

import (
	"load_balancer/algorithm"
	"net/http"
	"time"
)

func main () {
	
	/*
	I want to have a go server that has a list of other urls 
	(the actual http server to offload requests to)
	Simply start by just implmenting a round robin policy	rr.Serve()(w,r)	
	*/
	s := &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	rr := algorithm.NewRoundRobinLoadBalancer()
	
	/*
	Backend registering API, The backend sends a registering message with 
	the port number where it will be listening for requests coming from the
	load balancer. 
	*/
	http.HandleFunc("/register", rr.HandleBackendRegister)
	
	/* 
	Redirect the request to a backend server
	- Select one of the active backend servers from the pool. 
	- Make a new request for the backend.
	- Redirect the response of the backend to the client.
	*/
	http.HandleFunc("/lb", func(w http.ResponseWriter, r *http.Request) {
		rr.Serve()(w,r)	
	})
	s.ListenAndServe()
}
