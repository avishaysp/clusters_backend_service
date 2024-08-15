package CacheServiceClient

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
	configFile := "../service/CacheServiceRoot.yaml"
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

func TestSetGetAndDelete(t *testing.T) {
	c := NewCacheServiceClient(conf.RegistryAddresses, conf.Type)
	err := c.Set("key1", "value1")
	if err != nil {
		t.Fatalf("could not call Set: %v", err)
		return
	}
	r, err := c.Get("key1")
	if err != nil {
		t.Fatalf("could not call Get: %v", err)
		return
	}
	if r != "value1" {
		t.Fatalf("wrong value: received %s, expected value1", r)
		return
	}
	err = c.Delete("key1")
	if err != nil {
		t.Fatalf("could not call Delete: %v", err)
		return
	}
	r, err = c.Get("key1")
	if err != nil {
		t.Fatalf("could not call Get: %v", err)
		return
	}
	if r != "" {
		t.Fatalf("wrong value: received %s, expected '' ", r)
		return
	}
	t.Logf("Response: %v", r)
}

func TestIsAlive(t *testing.T) {
	c := NewCacheServiceClient(conf.RegistryAddresses, conf.Type)
	r, err := c.IsAlive()
	if err != nil {
		t.Fatalf("could not call IsAlive: %v", err)
		return
	}
	if !r {
		t.Fatalf("IsAlive returned false")
		return
	}
	t.Logf("Response: %v", r)
}
