package types

// Config represents the structure of the configuration file
type Config struct {
	Clusters []Cluster `yaml:"clusters"`
}

// Cluster represents a Kubernetes cluster configuration
type Cluster struct {
	Name         string       `yaml:"name"`
	Controlplane Controlplane `yaml:"controlplane"`
	Compute      Compute      `yaml:"compute"`
	Vip          Vip          `yaml:"vip"`
}

type Controlplane struct {
	Cluster              bool `yaml:"cluster"`
	Num                  int
	Core                 int    `yaml:"core"`
	Memory               int    `yaml:"memory"`
	BaseVmid             int    `yaml:"base_vmid"`
	AddressSansLastOctet string `yaml:"address_sans_last_octet"`
	LastOctet            int    `yaml:"last_octet"`
	GatewayLastOctet     int    `yaml:"gateway_last_octet"`
	Domain               string `yaml:"domain"`
	Nameserver           string `yaml:"nameserver"`
}

type Compute struct {
	Num                  int    `yaml:"num"`
	Core                 int    `yaml:"core"`
	Memory               int    `yaml:"memory"`
	BaseVmid             int    `yaml:"base_vmid"`
	AddressSansLastOctet string `yaml:"address_sans_last_octet"`
	LastOctet            int    `yaml:"last_octet"`
	GatewayLastOctet     int    `yaml:"gateway_last_octet"`
	Domain               string `yaml:"domain"`
	Nameserver           string `yaml:"nameserver"`
}

type Vip struct {
	Controlplane     bool   `yaml:"controlplane"`
	ControlplaneIp   string `yaml:"controlplane_ip"`
	ControlplaneFqdn string `yaml:"controlplane_fqdn"`
	Services         bool   `yaml:"services"`
}
