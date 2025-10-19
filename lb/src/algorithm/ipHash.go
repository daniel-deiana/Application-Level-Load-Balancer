// This file implements a stateful load balancing, that means that request coming from one user are always routed to the same backend 
package algorithm

import (
	"lb/datamodel"
	"net/http/httputil"
)

type ipHashingLoadBalancer struct {
	backends	map[int]*datamodel.BackendServer
	proxies		[]*httputil.ReverseProxy
}