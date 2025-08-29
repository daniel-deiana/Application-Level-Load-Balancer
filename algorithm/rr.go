package algorithm

import (
	"load_balancer/datamodel"
	"net/http/httputil"
	"net/http"
	"sync/atomic"
	"io"
	"fmt"
	"encoding/json"
)

/*	
- Selects the next server based on a simple round robin algorithm
*/

type RoundRobinLoadBalancer struct {
	servers map[int64]*datamodel.BackendServer
	index 	int64
	length  int64
}

func (rr *RoundRobinLoadBalancer) showState() {
	for index, server := range rr.servers {
		fmt.Printf("the server at index %d is %s",index, server.Host)
	}
}

func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
		return &RoundRobinLoadBalancer {
		servers : make(map[int64]*datamodel.BackendServer), 
		index 	: 0,
		length  : 0,
	}
}

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


