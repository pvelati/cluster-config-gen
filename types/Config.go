package types

// Config rappresenta la struttura del file di configurazione.
type Config struct {
	Clusters []Cluster `yaml:"clusters"`
}

// Cluster rappresenta una configurazione di cluster Kubernetes.
type Cluster struct {
	Name                       string `yaml:"name"`
	NumMaster                  int    `yaml:"num_master"`
	NumWorker                  int    `yaml:"num_worker"`
	VIP                        bool   `yaml:"vip"`
	MasterBaseVmid             int    `yaml:"master_base_vmid"`
	MasterAddressSansLastOctet string `yaml:"master_address_sans_last_octet"`
	MasterLastOctet            int    `yaml:"master_last_octet"`
	MasterGateway              int    `yaml:"master_gateway"`
	MasterDomain               string `yaml:"master_domain"`
	WorkerBaseVmid             int    `yaml:"worker_base_vmid"`
	WorkerAddressSansLastOctet string `yaml:"worker_address_sans_last_octet"`
	WorkerLastOctet            int    `yaml:"worker_last_octet"`
	WorkerGateway              int    `yaml:"worker_gateway"`
	WorkerDomain               string `yaml:"worker_domain"`
}
