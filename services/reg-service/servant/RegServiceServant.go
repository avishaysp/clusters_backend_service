package RegServiceServant

import (
	"fmt"
	"strings"
	"sync"
	"time"

	dht "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/servant/dht"
	"github.com/TAULargeScaleWorkshop/RLAD/utils"
)

// globals
var (
	is_first  bool
	chordNode *dht.Chord
	cacheMap  map[string]*NodeStatus // node address : NodeStatus
	mutex     sync.Mutex
)

type NodeStatus struct {
	FailCount int
	Alive     bool
}

func encodeStrings(lst []string) string {
	return strings.Join(lst, ",")
}

func decodeStrings(enc string) []string {
	return strings.Split(enc, ",")
}

func isInChord(key string) bool {
	keys, err := chordNode.GetAllKeys()
	if err != nil {
		utils.Logger.Fatalf("chordNode.GetAllKeys failed with error: %v", err)
	}

	// check if the service is in the keys list
	for _, item := range keys {
		if item == key {
			return true
		}
	}
	return false
}

// helper functions
func IsFirst() bool {
	return is_first
}

func InitServant(chord_name string) {
	utils.Logger.Printf("RegServiceServant::InitServant() called with %s", chord_name)
	var err error

	if chord_name == "root" {
		chordNode, err = dht.NewChord(chord_name, 1099)
		if err != nil {
			utils.Logger.Fatalf("could not create new chord: %v", err)
			return
		}
		utils.Logger.Printf("NewChord returned: %v", chordNode)
		cacheMap = make(map[string]*NodeStatus)
	} else {
		// join already existing "root" with a new chord_name
		chordNode, err = dht.JoinChord(chord_name, "root", 1099)
		if err != nil {
			utils.Logger.Fatalf("could not join chord: %v", err)
			return
		}
		utils.Logger.Printf("JoinChord returned: %v", chordNode)
	}

	is_first, err = chordNode.IsFirst()
	if err != nil {
		utils.Logger.Fatalf("could not call IsFirst: %v", err)
		return
	}
	utils.Logger.Printf("chordNode.IsFirst() result: %v", is_first)
}

// Registry API
func Register(service_name string, node_address string) {
	mutex.Lock()
	defer mutex.Unlock()

	var addresses []string // by default, empty list
	var err error

	// get service addresses
	if isInChord(service_name) {
		// get the current list
		enc, err := chordNode.Get(service_name)
		if err != nil {
			utils.Logger.Printf("chordNode.Get failed with error: %v", err)
		}
		addresses = decodeStrings(enc)
	}

	// checks if address already exists
	if len(addresses) > 0 {
		for _, address := range addresses {
			if address == node_address {
				utils.Logger.Printf("Address %s already exists for service %s\n", node_address, service_name)
				return
			}
		}
	}

	// add to list and set back to chord
	addresses = append(addresses, node_address)
	updated_enc := encodeStrings(addresses)
	err = chordNode.Set(service_name, updated_enc)
	if err != nil {
		utils.Logger.Printf("chordNode.Set failed with error: %v", err)
	}
	utils.Logger.Printf("Address %s added for service %s\n", node_address, service_name)
}

// note: assuming service_name is registered
func unregisterFromChord(service_name string, node_address string) {
	// get the current list
	enc, err := chordNode.Get(service_name)
	if err != nil {
		utils.Logger.Printf("chordNode.Get failed with error: %v", err)
	}
	lst := decodeStrings(enc)
	if len(lst) == 0 {
		utils.Logger.Printf("Service %s not found\n", service_name)
		return
	}

	for i, address := range lst {
		if address == node_address {
			// remove from list and set back to chord
			lst = append(lst[:i], lst[i+1:]...)
			if len(lst) == 0 {
				err = chordNode.Delete(service_name)
				if err != nil {
					utils.Logger.Printf("chordNode.Delete failed with error: %v", err)
				}
				return
			}
			utils.Logger.Printf("Address %s removed for service %s\n", node_address, service_name)
			updated_enc := encodeStrings(lst)
			err = chordNode.Set(service_name, updated_enc)
			if err != nil {
				utils.Logger.Printf("chordNode.Set failed with error: %v", err)
			}
			return
		}
	}
	utils.Logger.Printf("Address %s not found for service %s\n", node_address, service_name)
}

func Unregister(service_name string, node_address string) {
	mutex.Lock()
	defer mutex.Unlock()

	if !isInChord(service_name) {
		utils.Logger.Printf("Service %s not registered!", service_name)
	}

	unregisterFromChord(service_name, node_address)
}

func Discover(service_name string) ([]string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if !isInChord(service_name) {
		return nil, fmt.Errorf("service %s not registered", service_name)
	}

	// get the current list
	enc, err := chordNode.Get(service_name)
	if err != nil {
		utils.Logger.Printf("chordNode.Get failed with error: %v", err)
	}
	lst := decodeStrings(enc)
	if len(lst) == 0 {
		return lst, fmt.Errorf("service not found: %v", service_name)
	}
	return lst, nil
}

// Internal logic, health checking the nodes
func IsAliveCheck() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		mutex.Lock()
		utils.Logger.Printf("IsAliveCheck: called\n")

		// get all the services
		services, err := chordNode.GetAllKeys()
		if err != nil {
			utils.Logger.Fatalf("chordNode.GetAllKeys failed with error: %v", err)
		}

		for _, serviceName := range services {
			// get the current list
			enc, err := chordNode.Get(serviceName)
			if err != nil {
				utils.Logger.Printf("chordNode.Get failed with error: %v", err)
			}
			addresses := decodeStrings(enc)

			for _, address := range addresses {
				utils.Logger.Printf("IsAliveCheck: Service = %s, Node = %v\n", serviceName, address)
				var alive bool
				var err error
				switch serviceName {
				case "TestService": // grpc
					c := NewTestServiceClient(address, "")
					alive, err = c.IsAlive()
				case "TestServiceMQ": // mq
					c := NewTestServiceClient("", address)
					r, err := c.IsAliveAsync()
					if err != nil {
						utils.Logger.Printf("could not call IsAliveAsync: %v", err)
						continue
					}
					alive, err = r()
					//alive = true
				case "CacheService":
					c := NewCacheServiceClient(address)
					alive, err = c.IsAlive()
				default:
					utils.Logger.Printf("Unknown service name: %v", serviceName)
				}

				// create node status if doesn't exist
				_, ok := cacheMap[address]
				if !ok {
					cacheMap[address] = &NodeStatus{FailCount: 0, Alive: true}
				}

				// we assume that (!alive) iff (err != nil) in IsAlive implementation
				if !alive {
					utils.Logger.Printf("IsAliveCheck: Node %v is not alive: error = %v\n", address, err)
					cacheMap[address].FailCount++
					if cacheMap[address].FailCount >= 2 {
						utils.Logger.Printf("IsAliveCheck: marking node as not alive! %v\n", address)
						// mark the node to be unregistered later
						cacheMap[address].Alive = false
					}
				} else {
					utils.Logger.Printf("IsAliveCheck: Node %v is alive", address)
					cacheMap[address].FailCount = 0
					cacheMap[address].Alive = true
				}
			}

			// unregister (manually, not calling the API function) the "dead" nodes
			for i := len(addresses) - 1; i >= 0; i-- {
				if !cacheMap[addresses[i]].Alive {
					utils.Logger.Printf("Node %s is not alive, unregistering...\n", addresses[i])
					// remove from cache
					delete(cacheMap, addresses[i])

					// unregister
					unregisterFromChord(serviceName, addresses[i])
				}
			}
		}
		mutex.Unlock()
	}
}
