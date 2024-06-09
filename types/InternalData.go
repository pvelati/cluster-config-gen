package types

// InternalData rappresenta i dati interni generati dal processo di generazione di configurazioni.
type InternalData struct {
	Clusters []InternalDataCluster
}

// InternalDataCluster rappresenta i dati di un singolo cluster nell'InternalData.
type InternalDataCluster struct {
	Name    string
	Masters []InternalDataMaster
	Workers []InternalDataWorker

	AnsibleMasterGroup string
	AnsibleWorkerGroup string
}

// InternalDataMaster rappresenta i dati di un master node all'interno di un cluster.
type InternalDataMaster struct {
	Host string
	IP   string

	ProxmoxVMID          int
	ProxmoxVmName        string
	ProxmoxVmDescription string
	ProxmoxVmTags        []string

	TerraformResourceName string
}

// InternalDataWorker rappresenta i dati di un worker node all'interno di un cluster.
type InternalDataWorker struct {
	Host string
	IP   string

	ProxmoxVMID          int
	ProxmoxVmName        string
	ProxmoxVmDescription string
	ProxmoxVmTags        []string

	TerraformResourceName string
}
