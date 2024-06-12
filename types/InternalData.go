package types

// InternalData rappresenta i dati interni generati dal processo di generazione di configurazioni.
type InternalData struct {
	Clusters []InternalDataCluster
}

// InternalDataCluster rappresenta i dati di un singolo cluster nell'InternalData.
type InternalDataCluster struct {
	Name    string
	Masters []InternalDataNode
	Workers []InternalDataNode
	Ha      bool
	HaIp    string

	AnsibleMasterGroup string
	AnsibleWorkerGroup string
}

// InternalDataNode rappresenta i dati di un nodo generico all'interno di un cluster.
type InternalDataNode struct {
	Host    string
	Domain  string
	IP      string
	Gateway string

	ProxmoxVMID          int
	ProxmoxVmName        string
	ProxmoxVmDescription string
	ProxmoxVmTags        []string

	TerraformResourceName string
}
