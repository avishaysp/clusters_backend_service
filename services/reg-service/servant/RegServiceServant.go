package RegServiceServant

import (
	"fmt"
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

// helper functions
func IsFirst() bool {
	return is_first
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
func Register(service_name string, addresses NodeAddresses) {
	mutex.Lock()
	defer mutex.Unlock()

	var nodes []string // by default, empty list
	var err error

	// get service addresses
	if isInChord(service_name) {
		// get the current list
		enc, err := chordNode.Get(service_name)
		if err != nil {
			utils.Logger.Printf("chordNode.Get failed with error: %v", err)
		}
		nodes = decodeStrings(enc)
	}

	// encode addresses
	enc_node_address := encodeProtocols(addresses)

	// checks if (encoded) address already exists
	if len(nodes) > 0 {
		for _, node_address := range nodes {
			if node_address == enc_node_address {
				utils.Logger.Printf("Address %s already exists for service %s\n", enc_node_address, service_name)
				return
			}
		}
	}

	// add to list and set back to chord
	nodes = append(nodes, enc_node_address)
	updated_enc := encodeStrings(nodes)
	err = chordNode.Set(service_name, updated_enc)
	if err != nil {
		utils.Logger.Printf("chordNode.Set failed with error: %v", err)
	}
	utils.Logger.Printf("Address %s added for service %s\n", enc_node_address, service_name)
}

// note: assuming service_name is registered
func unregisterFromChord(service_name string, addresses NodeAddresses) {
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

	// encode address with protocol to search for it
	enc_node_address := encodeProtocols(addresses)

	for i, address := range lst {
		if address == enc_node_address {
			// remove from list and set back to chord
			lst = append(lst[:i], lst[i+1:]...)
			if len(lst) == 0 {
				err = chordNode.Delete(service_name)
				if err != nil {
					utils.Logger.Printf("chordNode.Delete failed with error: %v", err)
				}
				return
			}
			utils.Logger.Printf("Address %s removed for service %s\n", enc_node_address, service_name)
			updated_enc := encodeStrings(lst)
			err = chordNode.Set(service_name, updated_enc)
			if err != nil {
				utils.Logger.Printf("chordNode.Set failed with error: %v", err)
			}
			return
		}
	}
	utils.Logger.Printf("Address %s not found for service %s\n", enc_node_address, service_name)
}

func Unregister(service_name string, addresses NodeAddresses) {
	mutex.Lock()
	defer mutex.Unlock()

	if !isInChord(service_name) {
		utils.Logger.Printf("Service %s not registered!", service_name)
	}

	unregisterFromChord(service_name, addresses)
}

// returns a list of strings in the format of "$<protocol>$<address>...$...$..."
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

			for _, enc_address := range addresses {
				utils.Logger.Printf("IsAliveCheck: Service = %s, Node = %v\n", serviceName, enc_address)

				// decode the address with its protocol
				node_addresses := decodeProtocols(enc_address)

				var alive bool
				var err error
				switch serviceName {
				case "TestService":
					c := NewTestServiceClient(node_addresses["GRPC"], "")
					alive, err = c.IsAlive()
				case "CacheService":
					c := NewCacheServiceClient(node_addresses["GRPC"])
					alive, err = c.IsAlive()
				default:
					utils.Logger.Printf("Unknown service name: %v", serviceName)
				}

				// create node status if doesn't exist
				_, ok := cacheMap[enc_address]
				if !ok {
					cacheMap[enc_address] = &NodeStatus{FailCount: 0, Alive: true}
				}

				// we assume that (!alive) iff (err != nil) in IsAlive implementation
				if !alive {
					utils.Logger.Printf("IsAliveCheck: Node %v is not alive: error = %v\n", node_addresses["GRPC"], err)
					cacheMap[enc_address].FailCount++
					if cacheMap[enc_address].FailCount >= 2 {
						utils.Logger.Printf("IsAliveCheck: marking node as not alive! %v\n", node_addresses["GRPC"])
						// mark the node to be unregistered later
						cacheMap[enc_address].Alive = false
					}
				} else {
					utils.Logger.Printf("IsAliveCheck: Node %v is alive", enc_address)
					cacheMap[enc_address].FailCount = 0
					cacheMap[enc_address].Alive = true
				}
			}

			// unregister (manually, not calling the API function) the "dead" nodes
			for i := len(addresses) - 1; i >= 0; i-- {
				if !cacheMap[addresses[i]].Alive {
					utils.Logger.Printf("Node %s is not alive, unregistering...\n", addresses[i])
					// remove from cache
					delete(cacheMap, addresses[i])

					// unregister
					unregisterFromChord(serviceName, decodeProtocols(addresses[i]))
				}
			}
		}
		mutex.Unlock()
	}
}
