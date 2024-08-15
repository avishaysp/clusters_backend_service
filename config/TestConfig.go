package config

type TestConfig struct {
	BaseConfig
	RegistryAddresses []string `yaml:"registry_addresses"`
}
