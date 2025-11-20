package algorithm
import (
	"lb/helper"
	"lb/datamodel"
	"net/http/httputil"
	"net/http"
	"sync/atomic"
	"io"
	"fmt"
	"encoding/json"
	"time"
	"net"
)

/*	
- Selects the next server based on a simple round robin algorithm
*/
type RoundRobinLoadBalancer struct {
	servers 			map[string]*datamodel.BackendServer
	index 				int64
	length  			int64
	proxies 			[]*httputil.ReverseProxy
	tickerRefreshStates *time.Ticker
	stopMonitoring 		chan bool
}

func (rr *RoundRobinLoadBalancer) showState() {
	for host, server := range rr.servers {
		fmt.Printf("the server at index %s is %s\n",host, server.Host)
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
	strURL := "http://" + bs.Host + ":8081" + "/health"	
	
	client := http.Client{
    	Timeout: 2* time.Second,
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
		atomic.AddInt64(&rr.length,1)
	} else if (resp.StatusCode != 200 && bs.State == datamodel.StateConnected) {
		newState = datamodel.StateError 
		rr.proxies = helper.RemoveByValue(rr.proxies, bs.Proxy)
		atomic.AddInt64(&rr.length,-1)
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
		rr.addNewBackend(r.RemoteAddr)
		rr.showState()
}

func (rr *RoundRobinLoadBalancer) addNewBackend (Host string) {
	var protocols = []string{"http"}

	// NB: if you try to search a key in the dict it returns the value and the presence
	if _, isPresent := rr.servers[Host]; isPresent {
		fmt.Printf("The server is already registered in the load balancer")
		return
	}

	Host, _, _ = net.SplitHostPort(Host)
	// if (err != nil) {
	// 	fmt.Printf("some error in splitting the host:port format from backend")
	// }

	myRewrite := func(pr *httputil.ProxyRequest) {
		// Print all incoming client request headers
		fmt.Printf("the client URL.Host is %s", pr.In.Host)
		pr.Out.URL.Scheme = "http"
		pr.Out.URL.Host = Host + ":8081"
		fmt.Printf("the backend server im forwarding the request is %s", pr.Out.URL.Host)
		pr.Out.Header.Set("X-Forwarded-Host", pr.In.RemoteAddr)
	}

	fmt.Printf("the key im using to register the backend is %s", Host)

	newBs := datamodel.NewBackendServer(Host,protocols, myRewrite)
	rr.servers[Host] = newBs	
	
	rr.proxies = append(rr.proxies, newBs.Proxy)
	// add a new proxy to the list of current backends
	atomic.AddInt64(&rr.length,1)
}

func (rr *RoundRobinLoadBalancer) getNext() *httputil.ReverseProxy {
	atomic.AddInt64(&rr.index,1)
	return rr.proxies[rr.index%rr.length]
}

func (rr *RoundRobinLoadBalancer) Serve() func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		rr.getNext().ServeHTTP(w,r)
	} 
}