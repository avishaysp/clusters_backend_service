package RegServiceClient

import (
	"log"
	"os"
	"testing"

	"github.com/TAULargeScaleWorkshop/RLAD/config"
	"gopkg.in/yaml.v2"
)

// global config
var conf config.TestConfig

// TestMain is the entry point for testing
func TestMain(m *testing.M) {
	// Load the configuration
	configFile := "../../test-service/service/TestService.yaml"
	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
		os.Exit(2)
	}

	// get base config
	var baseConfig config.BaseConfig
	err = yaml.Unmarshal(configData, &baseConfig) // parses YAML
	if err != nil {
		log.Fatalf("error unmarshaling BaseConfig data: %v", err)
	}

	// get TestService config
	err = yaml.Unmarshal(configData, &conf) // parses YAML
	if err != nil {
		log.Fatalf("error unmarshaling TestConfig data: %v", err)
	}
	conf.BaseConfig = baseConfig
	log.Printf("loaded config %s", configData)

	// Run the tests
	code := m.Run()

	// Exit with the test run's exit code
	os.Exit(code)
}

func TestRegisterUnregister(t *testing.T) {
	c := NewRegServiceClient(conf.RegistryAddresses)
	err := c.Register("test1", "node1")
	if err != nil {
		t.Fatalf("could not call Register: %v", err)
		return
	}
	nodes, err := c.Discover("test1")
	if err != nil {
		t.Fatalf("could not call Discover: %v", err)
		return
	}
	if len(nodes) != 1 && nodes[0] != "node1" {
		t.Fatalf("wrong value for discover after register: received %s, expected node1", nodes[0])
		return
	}
	err = c.Unregister("test1", "node1")
	if err != nil {
		t.Fatalf("could not call Unregister: %v", err)
		return
	}
	// we expect nodes to be nil
	nodes, _ = c.Discover("test1")
	if len(nodes) != 0 {
		t.Fatalf("wrong value for discover after register: expected no addresses for this service")
		return
	}
	t.Logf("Successful Register and Unregister")
}

func TestDifferentRegNodes(t *testing.T) {
	first_node := []string{"127.0.0.1:8502"}
	c := NewRegServiceClient(first_node)
	err := c.Register("test1", "node1")
	if err != nil {
		t.Fatalf("could not call Register: %v", err)
		return
	}

	second_node := []string{"127.0.0.1:8503"}
	c2 := NewRegServiceClient(second_node)
	nodes, _ := c2.Discover("test1")
	if len(nodes) == 0 {
		t.Fatalf("expected 1 node")
		return
	}
	log.Printf("Nodes = %s", nodes[0])
}
