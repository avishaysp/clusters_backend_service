package CacheServiceServant

import (
	dht "github.com/TAULargeScaleWorkshop/RLAD/services/cache-service/servant/dht"
	"github.com/TAULargeScaleWorkshop/RLAD/utils"
)

// globals
var (
	is_first  bool
	chordNode *dht.Chord
)

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
	utils.Logger.Printf("IsFirst() called, result: %v", is_first)
	return is_first
}

func InitServant(chord_name string) {
	utils.Logger.Printf("CacheServiceServant::InitServant() called with %v", chord_name)
	var err error

	if chord_name == "root" {
		chordNode, err = dht.NewChord(chord_name, 2099)
		if err != nil {
			utils.Logger.Fatalf("could not create new chord: %v", err)
			return
		}
		utils.Logger.Printf("NewChord returned: %v", chordNode)
	} else {
		// join already existing "root" with a new chord_name
		chordNode, err = dht.JoinChord(chord_name, "root", 2099)
		if err != nil {
			utils.Logger.Fatalf("could not join chord: %v", err)
			return
		}
		utils.Logger.Printf("JoinChord returned: %v", chordNode)
	}
	// TODO: consider removing later
	// check
	is_first, err = chordNode.IsFirst()
	if err != nil {
		utils.Logger.Fatalf("could not call IsFirst: %v", err)
		return
	}
	utils.Logger.Printf("chordNode.IsFirst() result: %v", is_first)
}

func Set(key string, value string) error {
	err := chordNode.Set(key, value)
	if err != nil {
		utils.Logger.Printf("chordNode.Set failed with error: %v", err)
		return err
	}
	utils.Logger.Printf("Value %s added for key %s\n", value, key)
	return nil
}

func Get(key string) (string, error) {
	// returns the value (or "" if not found), and a boolean indicating whether the key was found in the chord
	var err error
	var value string
	if isInChord(key) {
		// get the current list
		value, err = chordNode.Get(key)
		if err != nil {
			utils.Logger.Printf("chordNode.Get failed with error: %v", err)
		}
	}
	return value, err
}

func Delete(key string) error {
	err := chordNode.Delete(key)
	if err != nil {
		utils.Logger.Printf("chordNode.Delete failed with error: %v", err)
	}
	return err
}

func IsAlive() bool {
	return true
}
