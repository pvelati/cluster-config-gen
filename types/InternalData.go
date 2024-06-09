package types

type InternalData struct {
	Clusters []InternalDataCluster
}

type InternalDataCluster struct {
	Name    string
	Masters []InternalDataMaster
	Workers []InternalDataWorker

	AnsibleMasterGroup string
	AnsibleWorkerGroup string
}

type InternalDataMaster struct {
	Host string
	IP   string

	ProxmoxVMID          int
	ProxmoxVmName        string
	ProxmoxVmDescription string
	ProxmoxVmTags        []string

	TerraformResourceName string
}

type InternalDataWorker struct {
	Host string
	IP   string

	ProxmoxVMID          int
	ProxmoxVmName        string
	ProxmoxVmDescription string
	ProxmoxVmTags        []string

	TerraformResourceName string
}
