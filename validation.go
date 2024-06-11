package main

import (
	"fmt"

	"github.com/pvelati/cluster-config-gen/types"
)

// validateClusters checks the clusters for any configuration issues.
func validateClusters(clusters []types.Cluster) error {
	nameSet := make(map[string]struct{})

	for i := range clusters {
		cluster := &clusters[i]

		if cluster.Name == "" {
			return fmt.Errorf("%s - name cannot be empty", cluster.Name)
		}
		if _, exists := nameSet[cluster.Name]; exists {
			return fmt.Errorf("duplicate cluster name: %s", cluster.Name)
		}
		nameSet[cluster.Name] = struct{}{}

		if cluster.NumMaster != 1 && cluster.NumMaster != 3 {
			return fmt.Errorf("%s - the number of masters must be 1 or 3", cluster.Name)
		}
		if cluster.NumWorker < 1 {
			return fmt.Errorf("%s - the cluster must have at least 1 worker", cluster.Name)
		}
		if cluster.MasterBaseVmid <= 0 {
			return fmt.Errorf("%s - master_base_vmid should be a positive integer", cluster.Name)
		}
		if cluster.MasterLastOctet <= 0 || cluster.MasterLastOctet >= 256 {
			return fmt.Errorf("%s - master_last_octet should be between 1 and 255", cluster.Name)
		}
		if cluster.MasterGateway <= 0 || cluster.MasterGateway >= 256 {
			return fmt.Errorf("%s - master_gateway should be between 1 and 255", cluster.Name)
		}
		if cluster.MasterDomain == "" {
			return fmt.Errorf("%s - master_domain cannot be empty", cluster.Name)
		}

		// Worker fallback to Master values if empty
		if cluster.WorkerBaseVmid == 0 {
			cluster.WorkerBaseVmid = cluster.MasterBaseVmid
		}
		if cluster.WorkerAddressSansLastOctet == "" {
			cluster.WorkerAddressSansLastOctet = cluster.MasterAddressSansLastOctet
		}
		if cluster.WorkerLastOctet == 0 {
			cluster.WorkerLastOctet = cluster.MasterLastOctet
		}
		if cluster.WorkerGateway == 0 {
			cluster.WorkerGateway = cluster.MasterGateway
		}
		if cluster.WorkerDomain == "" {
			cluster.WorkerDomain = cluster.MasterDomain
		}

		if cluster.WorkerBaseVmid <= 0 {
			return fmt.Errorf("worker_base_vmid should be a positive integer")
		}
		if cluster.WorkerLastOctet <= 0 || cluster.WorkerLastOctet >= 256 {
			return fmt.Errorf("worker_last_octet should be between 1 and 255")
		}
		if cluster.WorkerGateway <= 0 || cluster.WorkerGateway >= 256 {
			return fmt.Errorf("worker_gateway should be between 1 and 255")
		}
	}
	return nil
}
