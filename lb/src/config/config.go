package config

import (
	"fmt"
    "log"
    "github.com/fsnotify/fsnotify"
)

// object that holds configuration data
type LoadBalancerConfig struct {
	MaxBackends			int
	APIRateLimit		float64
	HealthCheckInterval	int
	BackendFailTimeout	int
}

func NewLoadBalancerConfig() *LoadBalancerConfig {
	lbConfig := &LoadBalancerConfig{
		MaxBackends:        10,
		APIRateLimit:       3.0,
		HealthCheckInterval: 2,
		BackendFailTimeout:  2,
	}
	return lbConfig
}

// class that handles initialization and update of config object
type LoadBalancerConfigManager struct{
	watcher*	fsnotify.Watcher
	config		LoadBalancerConfig
}

// creates a new config manager that will watch for update in the directory pointed by path 
func NewLoadBalancerConfigManager() *LoadBalancerConfigManager {
	manager := &LoadBalancerConfigManager{}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("watcher creation has returned error")
	}
	manager.watcher = watcher
	// should initizalize configuration object here
	return manager
}

func (manager *LoadBalancerConfigManager) StartWatchingConfigUpdates() {
	
	// start the goroutine that watches the update on config file
	go func() {
		for {
			select {
				case event, ok := <-manager.watcher.Events:
					if !ok {
						return 
					}
					log.Println("event: ", event)
					fmt.Printf("updated a file in the directory")
					if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
						// update config by reading modified file 
						log.Println("a file has been modified/added in the watched directory")
					}
				case err, ok := <-manager.watcher.Errors:
					if !ok {
						log.Println("The notify service has returned an error", err)
						return
					}
			}
		}
	}()

	err := manager.watcher.Add("/app")

	if err != nil {
		log.Fatal(err)
	}
}
