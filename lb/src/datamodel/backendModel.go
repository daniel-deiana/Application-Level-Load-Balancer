package datamodel

import (
	"net/http/httputil"
	"net/url"
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

type BackendState int 

/*
 define constant of type backendState with incremental values
 Each constant is one of the values of our enum
*/

const (
	StateIdle BackendState = iota
	StateConnected
	StateError
	StateSleeping
)

var stateMapping = map[BackendState]string {
	StateIdle 			: "idle",
	StateConnected 		: "connected",
	StateError			: "error",
	StateSleeping		: "sleep",
}

/*
	We are implementing the String() method that implements the Stringer 
	interface, this is used to print our enum returning a string
*/
func (bs BackendState) String() string {
	return stateMapping[bs] 
}

type BackendServer struct {
	State 		BackendState
	Host 		string
	Protocols 	[]string
	Proxy    	*httputil.ReverseProxy
}

func NewBackendServer(newHost string, newProtocols []string, rewriteHandler func(pr* httputil.ProxyRequest)) (*BackendServer){
	h,_ := url.Parse("http://"+newHost)
	var newProxy = httputil.NewSingleHostReverseProxy(h)
	// update custom write 
	newProxy.Director = nil
	newProxy.Rewrite = rewriteHandler
	var newState = StateConnected
	return &BackendServer {
		State : newState , 
		Host 		: newHost, 
		Protocols 	: newProtocols,
		Proxy  		: newProxy,
	}
}

