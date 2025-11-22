package main;

import (
	"lb/config"
	"lb/algorithm"
	"net/http"
	"time"
)

// Load Balancer interface defines a serve function common for all load balancers 
type ILoadBalancer interface{
	Serve() func(w http.ResponseWriter, r *http.Request)
}

func LoadBalancerServeEndpoint(endpoint string, lb ILoadBalancer) {
	lbEndpointHandler := lb.Serve()	
	http.HandleFunc(endpoint, lbEndpointHandler)
}

func main () {

	/*
		Start the configuration manager and get configuration of the load balancer
	*/
	lbConfManager := config.NewLoadBalancerConfigManager()

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
	lb := algorithm.NewRoundRobinLoadBalancer()
	
	lbConfManager.StartWatchingConfigUpdates()
	
	/*
	Backend registering API, The backend sends a registering message with 
	the port number where it will be listening for requests coming from the
	load balancer. 
	*/
	http.HandleFunc("/register", lb.HandleBackendRegister)
	
	/* 
	Redirect the request to a backend server
	- Select one of the active backend servers from the pool. 
	- Make a new request for the backend.
	- Redirect the response of the backend to the client.
	*/

	LoadBalancerServeEndpoint("/lb", lb)

	s.ListenAndServe()
}
