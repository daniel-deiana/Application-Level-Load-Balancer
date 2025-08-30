package datamodel

import (
	"net/http/httputil"
	"net/url"
	"net/http"
	"time"
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

type BackendServer struct {
	State 		backendState
	Host 		string
	Protocols 	[]string
	Proxy    	*httputil.ReverseProxy
}

func NewBackendServer(newHost string, newProtocols []string) (*BackendServer){
	h,_ := url.Parse("http://"+newHost)
	var newProxy = httputil.NewSingleHostReverseProxy(h)	
	var newState = StateConnected
	return &BackendServer {
		State : newState , 
		Host 		: newHost, 
		Protocols 	: newProtocols,
		Proxy  		: newProxy,
	}
}

/*
	Make a req to a backendserver and evaluate the new state depending on
	the response
*/
func (bs* BackendServer) UpdateState() {
	strURL := "http://" + bs.Host + "/health"	
	
	client := http.Client{
    	Timeout: 5 * time.Second,
	}
	resp, err := client.Get(strURL)
	if (err != nil) {
		bs.State = StateError
		return
	}	

	defer resp.Body.Close()

	if (resp.StatusCode == 200 && bs.State != StateConnected) {
		bs.State = StateConnected 
	} else if (resp.StatusCode != 200 && bs.State == StateConnected) {
		bs.State = StateError 
	}
}

