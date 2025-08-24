package algorithm

import (
	"load_balancer/datamodel"
)

/*	
	- Selects the next server based on a simple round robin algorithm.
*/
func GetNextServerRR(servers map[int]datamodel.BackendServer, current int) (serverHost string, new_current int) {
	server := servers[current]
	new_current = (current + 1) % len(servers)
	return server.Host, new_current
}