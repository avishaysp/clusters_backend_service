package TestServiceClient

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
	configFile := "../service/TestService.yaml"
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

func TestHelloWorld(t *testing.T) {
	c := NewTestServiceClient(conf.RegistryAddresses, conf.Type)
	r, err := c.HelloWorld()
	if err != nil {
		t.Fatalf("could not call HelloWorld: %v", err)
		return
	}
	t.Logf("Response: %v", r)
}

func TestHelloWorldAsync(t *testing.T) {
	c := NewTestServiceClient(conf.RegistryAddresses, conf.Type)
	r, err := c.HelloWorldAsync()
	if err != nil {
		t.Fatalf("could not call HelloWorldAsync: %v", err)
		return
	}
	res, err := r()
	if err != nil {
		t.Fatalf("HelloWorldAsync returned error : %v", err)
		return
	}
	t.Logf("Response: %v", res)
}

func TestHelloToUser(t *testing.T) {
	c := NewTestServiceClient(conf.RegistryAddresses, conf.Type)
	r, err := c.HelloToUser("Zvi")
	if err != nil {
		t.Fatalf("could not call HelloToUser: %v", err)
		return
	}
	t.Logf("Response: %v", r)
}

func TestHelloToUserAsync(t *testing.T) {
	c := NewTestServiceClient(conf.RegistryAddresses, conf.Type)
	r, err := c.HelloToUserAsync("ZviAsync")
	if err != nil {
		t.Fatalf("could not call HelloToUserAsync: %v", err)
		return
	}
	res, err := r()
	if err != nil {
		t.Fatalf("HelloToUserAsync returned error : %v", err)
		return
	}
	t.Logf("Response: %v", res)
}

func TestStoreAndGet(t *testing.T) {
	c := NewTestServiceClient(conf.RegistryAddresses, conf.Type)
	// store
	err := c.Store("key1", "value1")
	if err != nil {
		t.Fatalf("could not call Store: %v", err)
		return
	}
	// get
	r, err := c.Get("key1")
	if err != nil {
		t.Fatalf("could not call Get: %v", err)
		return
	}
	// check
	if r != "value1" {
		t.Fatalf("wrong value: received %s, expected value1", r)
		return
	}
	t.Logf("Response: %v", r)
}

func TestStoreAndGetAsync(t *testing.T) {
	c := NewTestServiceClient(conf.RegistryAddresses, conf.Type)

	// store
	r_store, err := c.StoreAsync("key1", "value1")
	if err != nil {
		t.Fatalf("could not call StoreAsync: %v", err)
		return
	}
	err = r_store()
	if err != nil {
		t.Fatalf("StoreAsync returned error : %v", err)
		return
	}
	// get
	r_get, err := c.GetAsync("key1")
	if err != nil {
		t.Fatalf("could not call GetAsync: %v", err)
		return
	}
	res, err := r_get()
	if err != nil {
		t.Fatalf("GetAsync returned error : %v", err)
		return
	}

	// check
	if res != "value1" {
		t.Fatalf("wrong value: received %s, expected value1", res)
		return
	}
	t.Logf("Response: %v", res)

}

// func TestWaitAndRand(t *testing.T) {
// 	c := NewTestServiceClient(conf.RegistryAddresses, conf.Type)
// 	resPromise, err := c.WaitAndRand(3)
// 	if err != nil {
// 		t.Fatalf("Calling WaitAndRand failed: %v", err)
// 		return
// 	}
// 	res, err := resPromise()
// 	if err != nil {
// 		t.Fatalf("WaitAndRand failed: %v", err)
// 		return
// 	}
// 	t.Logf("Returned random number: %v\n", res)
// }

func TestIsAlive(t *testing.T) {
	c := NewTestServiceClient(conf.RegistryAddresses, conf.Type)
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

func TestIsAliveAsync(t *testing.T) {
	c := NewTestServiceClient(conf.RegistryAddresses, conf.Type)
	r, err := c.IsAliveAsync()
	if err != nil {
		t.Fatalf("could not call IsAliveAsync: %v", err)
		return
	}
	res, err := r()
	if err != nil {
		t.Fatalf("IsAliveAsync returned error : %v", err)
		return
	}
	// check
	if !res {
		t.Fatalf("IsAlive returned false")
		return
	}
	t.Logf("Response: %v", res)
}

func TestExtractLinksFromURL(t *testing.T) {
	c := NewTestServiceClient(conf.RegistryAddresses, conf.Type)

	url := "https://www.microsoft.com"
	links, err := c.ExtractLinksFromURL(url, 1)
	if err != nil {
		t.Fatalf("ExtractLinksFromURL failed with error: %v", err)
	}

	// make sure you got some links
	if len(links) == 0 {
		t.Fatalf("ExtractLinksFromURL returned no links")
	}

	// print the links
	t.Logf("links: %v\n", links)
}

func TestExtractLinksFromURLAsync(t *testing.T) {
	c := NewTestServiceClient(conf.RegistryAddresses, conf.Type)
	url := "https://www.microsoft.com"
	r, err := c.ExtractLinksFromURLAsync(url, 1)
	if err != nil {
		t.Fatalf("could not call ExtractLinksFromURLAsync: %v", err)
		return
	}
	links, err := r()
	if err != nil {
		t.Fatalf("ExtractLinksFromURLAsync returned error : %v", err)
		return
	}

	// make sure you got some links
	if len(links) == 0 {
		t.Fatalf("ExtractLinksFromURL returned no links")
	}

	// print the links
	t.Logf("links: %v\n", links)
}
