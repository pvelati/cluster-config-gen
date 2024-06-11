package main

import (
	"fmt"
	"net"

	"github.com/pvelati/cluster-config-gen/types"
)

// calculateIpRange calculates the range of IPs given a base IP, starting last octet, and count of IPs.
func calculateIpRange(baseIP string, startOctet, count int) (types.IpRange, error) {
	if net.ParseIP(fmt.Sprintf("%s.1", baseIP)) == nil {
		return types.IpRange{}, fmt.Errorf("invalid base IP: %s", baseIP)
	}

	start := net.ParseIP(fmt.Sprintf("%s.%d", baseIP, startOctet))
	if start == nil {
		return types.IpRange{}, fmt.Errorf("invalid start IP from base: %s and octet: %d", baseIP, startOctet)
	}

	end := net.ParseIP(fmt.Sprintf("%s.%d", baseIP, startOctet+count-1))
	if end == nil {
		return types.IpRange{}, fmt.Errorf("invalid end IP from base: %s and octet: %d", baseIP, startOctet+count-1)
	}

	return types.IpRange{Start: start, End: end}, nil
}

// ipLessThan compares two IP addresses and returns true if ip1 is less than ip2.
func ipLessThan(ip1, ip2 net.IP) bool {
	for i := 0; i < len(ip1); i++ {
		if ip1[i] < ip2[i] {
			return true
		}
		if ip1[i] > ip2[i] {
			return false
		}
	}
	return false
}

// overlaps checks if two IP ranges overlap.
func overlaps(r1, r2 types.IpRange) bool {
	return !(ipLessThan(r1.End, r2.Start) || ipLessThan(r2.End, r1.Start))
}

// validateClusters checks the clusters for any configuration issues.
func validateClusters(clusters []types.Cluster) error {
	nameSet := make(map[string]struct{})
	var IpRanges []types.IpRange

	for i := range clusters {
		cluster := &clusters[i]

		if cluster.Name == "" {
			return fmt.Errorf("%s - name cannot be empty", cluster.Name)
		}
		if _, exists := nameSet[cluster.Name]; exists {
			return fmt.Errorf("duplicate cluster name: %s", cluster.Name)
		}
		nameSet[cluster.Name] = struct{}{}

		if cluster.NumWorker < 1 {
			return fmt.Errorf("%s - the cluster must have at least 1 worker", cluster.Name)
		}
		if cluster.MasterBaseVmid <= 0 {
			return fmt.Errorf("%s - master_base_vmid should be a positive integer", cluster.Name)
		}
		if cluster.MasterLastOctet <= 0 || cluster.MasterLastOctet >= 256 {
			return fmt.Errorf("%s - master_last_octet should be between 1 and 255", cluster.Name)
		}
		if cluster.MasterGatewayLastOctet <= 0 || cluster.MasterGatewayLastOctet >= 256 {
			return fmt.Errorf("%s - master_gateway_last_octet should be between 1 and 255", cluster.Name)
		}
		if cluster.MasterDomain == "" {
			return fmt.Errorf("%s - master_domain cannot be empty", cluster.Name)
		}
		if cluster.WorkerBaseVmid < 0 {
			return fmt.Errorf("%s - worker_base_vmid should be a positive integer", cluster.Name)
		}
		if cluster.WorkerLastOctet < 0 || cluster.WorkerLastOctet >= 256 {
			return fmt.Errorf("%s - worker_last_octet should be between 1 and 255", cluster.Name)
		}
		if cluster.WorkerGatewayLastOctet < 0 || cluster.WorkerGatewayLastOctet >= 256 {
			return fmt.Errorf("%s - worker_gateway_last_octet should be between 1 and 255", cluster.Name)
		}

		// Calculate and check IP ranges for masters
		masterRange, err := calculateIpRange(cluster.MasterAddressSansLastOctet, cluster.MasterLastOctet, cluster.NumMaster)
		if err != nil {
			return fmt.Errorf("error calculating master IP range for cluster %s: %v", cluster.Name, err)
		}

		// Calculate and check IP ranges for workers
		workerRange, err := calculateIpRange(cluster.WorkerAddressSansLastOctet, cluster.WorkerLastOctet, cluster.NumWorker)
		if err != nil {
			return fmt.Errorf("error calculating worker IP range for cluster %s: %v", cluster.Name, err)
		}

		// Check for IP range overlaps within the same cluster
		if overlaps(masterRange, workerRange) {
			return fmt.Errorf("IP range overlap detected within cluster %s between master and worker nodes", cluster.Name)
		}

		// Check for IP range overlaps across different clusters
		for _, existingRange := range IpRanges {
			if overlaps(masterRange, existingRange) || overlaps(workerRange, existingRange) {
				return fmt.Errorf("IP range overlap detected in cluster %s", cluster.Name)
			}
		}

		IpRanges = append(IpRanges, masterRange, workerRange)
	}
	return nil
}
