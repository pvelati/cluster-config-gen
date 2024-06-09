package types

// Config rappresenta la struttura del file di configurazione.
type Config struct {
	Clusters []Cluster `yaml:"clusters"`
}

// Cluster rappresenta una configurazione di cluster Kubernetes.
type Cluster struct {
	Name            string `yaml:"name"`
	NumMaster       int    `yaml:"num_master"`
	NumWorker       int    `yaml:"num_worker"`
	VIP             bool   `yaml:"vip"`
	MasterLastOctet int    `yaml:"master_last_octet"`
	WorkerLastOctet int    `yaml:"worker_last_octet"`
}
