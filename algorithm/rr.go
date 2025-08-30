package algorithm

import (
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
	servers 			map[int64]*datamodel.BackendServer
	index 				int64
	length  			int64
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
		servers : make(map[int64]*datamodel.BackendServer), 
		index 	: 0,
		length  : 0,
	}
	rrLoadBalancer.startMonitoringBackends()
	return &rrLoadBalancer
}

func (rr *RoundRobinLoadBalancer) startMonitoringBackends() {
	

	rr.tickerRefreshStates = time.NewTicker(1000 * time.Millisecond)
    rr.stopMonitoring = make(chan bool)
    go rr.HandleRefreshBackends()
}


func (rr *RoundRobinLoadBalancer) HandleRefreshBackends() {
    for {
    	// The select is a primitive that listens for go Channels 
        select {
        case <- rr.stopMonitoring:
            fmt.Println("Stopped monitoring backends\n")
            return   
        case t := <- rr.tickerRefreshStates.C:
			fmt.Println("Checking backend states at %v\n", t)
			rr.refreshBackendStates()
			rr.showState()     	
        }
    }
    rr.showState()
}

/*
	For every registerd server we check it's state by making a req 
	and update accordingly 
*/
func (rr *RoundRobinLoadBalancer) refreshBackendStates() {
	for _, server := range rr.servers {
		server.UpdateState()
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

func (rr *RoundRobinLoadBalancer) addNewBackend (Host string) {
	var protocols = []string{"http"}
	rr.servers[rr.length] = datamodel.NewBackendServer(Host,protocols)
	atomic.AddInt64(&rr.length,1)
}

func (rr *RoundRobinLoadBalancer) getNext() *httputil.ReverseProxy {
	atomic.AddInt64(&rr.index,1)
	return rr.servers[rr.index%rr.length].Proxy
}

func (rr *RoundRobinLoadBalancer) Serve() func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		rr.getNext().ServeHTTP(w,r)
	} 
}