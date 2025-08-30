package algorithm

import (
	"load_balancer/helper"
	"load_balancer/datamodel"
	"net/http/httputil"
	"net/http"
	"sync/atomic"
	"io"
	"fmt"
	"encoding/json"
	"time"
)

/*	
- Selects the next server based on a simple round robin algorithm
*/

type RoundRobinLoadBalancer struct {

	// TODO: Add a separate array for proxies so that the rr can loop on that 
	// This should be done to facilitate the rejoining of servers
	// currently using the integer as a key blocks a fast checking for duplicates 
	servers 			map[string]*datamodel.BackendServer
	index 				int64
	length  			int64
	proxies 			[]*httputil.ReverseProxy
	tickerRefreshStates *time.Ticker
	stopMonitoring 		chan bool
}

func (rr *RoundRobinLoadBalancer) showState() {
	for index, server := range rr.servers {
		fmt.Printf("the server at index %d is %s\n",index, server.Host)
		fmt.Printf("state : %s\n", server.State)
	}
}

func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
	var rrLoadBalancer = RoundRobinLoadBalancer {
		servers : make(map[string]*datamodel.BackendServer), 
		index 	: 0,
		length  : 0,
	}
	rrLoadBalancer.startMonitoringBackends()
	return &rrLoadBalancer
}

func (rr *RoundRobinLoadBalancer) startMonitoringBackends() {
	

	rr.tickerRefreshStates = time.NewTicker(1000 * time.Millisecond)
    rr.stopMonitoring = make(chan bool)
    go rr.HandleMonitorBackends()
}


func (rr *RoundRobinLoadBalancer) HandleMonitorBackends() {
    for {
    	// The select is a primitive that listens for go Channels 
        select {
        case <- rr.stopMonitoring:
            fmt.Println("Stopped monitoring backends\n")
            return   
        case t := <- rr.tickerRefreshStates.C:
			fmt.Println("Checking backend states at %v\n", t)
			rr.updateBackendStates()
			rr.showState()     	
        }
    }
}

/*
	For every registerd server we check it's state by making a req 
	and update accordingly 
*/
func (rr *RoundRobinLoadBalancer) updateBackendStates() {
	for _, server := range rr.servers {
		rr.updateBackendState(server)
	}	
}

/*
	Make a req to a backendserver and evaluate the new state depending on
	the response
*/
func (rr* RoundRobinLoadBalancer) updateBackendState(bs *datamodel.BackendServer) {
	strURL := "http://" + bs.Host + "/health"	
	
	client := http.Client{
    	Timeout: 5 * time.Second,
	}
	resp, err := client.Get(strURL)
		
	var newState datamodel.BackendState = datamodel.StateConnected

	defer func () {
		bs.State = newState
	}()

	if (err != nil) {
		newState = datamodel.StateError 
		return
	}	

	defer resp.Body.Close()

	if (resp.StatusCode == 200 && bs.State != datamodel.StateConnected) {
		newState = datamodel.StateConnected 
		rr.proxies = append(rr.proxies, bs.Proxy)
	} else if (resp.StatusCode != 200 && bs.State == datamodel.StateConnected) {
		newState = datamodel.StateError 
		rr.proxies = helper.RemoveByValue(rr.proxies, bs.Proxy)
	}

}

/*
	Handler for the registering API that registers a new server on the system
*/
func (rr *RoundRobinLoadBalancer) HandleBackendRegister(w http.ResponseWriter, r *http.Request) () {
		buf, _ := io.ReadAll(r.Body) 
		fmt.Printf("the buf contains %s\n", buf)
		var backendMsg datamodel.BackendRegMessage		
		json.Unmarshal(buf,&backendMsg)
		rr.addNewBackend(backendMsg.Host)
		rr.showState()
}


// todo refact 
func (rr *RoundRobinLoadBalancer) addNewBackend (Host string) {
	var protocols = []string{"http"}
	// add to backends map
	newBs := datamodel.NewBackendServer(Host,protocols)
	rr.servers[Host] = newBs	
	// add a new proxy to the list of current backends
	rr.proxies = append(rr.proxies, newBs.Proxy)
	atomic.AddInt64(&rr.length,1)
}


// todo refact 
func (rr *RoundRobinLoadBalancer) getNext() *httputil.ReverseProxy {
	atomic.AddInt64(&rr.index,1)
	return rr.proxies[rr.index%rr.length]
}

func (rr *RoundRobinLoadBalancer) Serve() func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		rr.getNext().ServeHTTP(w,r)
	} 
}