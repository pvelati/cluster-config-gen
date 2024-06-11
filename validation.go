package main

import (
	"fmt"
	"net"

	"github.com/pvelati/cluster-config-gen/types"
)

// ipRange represents a range of IPs.
type ipRange struct {
	start net.IP
	end   net.IP
}

// calculateIPRange calculates the range of IPs given a base IP, starting last octet, and count of IPs.
func calculateIPRange(baseIP string, startOctet, count int) (ipRange, error) {
	if net.ParseIP(fmt.Sprintf("%s.1", baseIP)) == nil {
		return ipRange{}, fmt.Errorf("invalid base IP: %s", baseIP)
	}

	start := net.ParseIP(fmt.Sprintf("%s.%d", baseIP, startOctet))
	if start == nil {
		return ipRange{}, fmt.Errorf("invalid start IP from base: %s and octet: %d", baseIP, startOctet)
	}

	end := net.ParseIP(fmt.Sprintf("%s.%d", baseIP, startOctet+count-1))
	if end == nil {
		return ipRange{}, fmt.Errorf("invalid end IP from base: %s and octet: %d", baseIP, startOctet+count-1)
	}

	return ipRange{start: start, end: end}, nil
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
func overlaps(r1, r2 ipRange) bool {
	return !(ipLessThan(r1.end, r2.start) || ipLessThan(r2.end, r1.start))
}

// validateClusters checks the clusters for any configuration issues.
func validateClusters(clusters []types.Cluster) error {
	nameSet := make(map[string]struct{})
	ipRanges := []ipRange{}

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
		// Check if worker last octet needs to be adjusted based on WorkerAddressSansLastOctet
		if cluster.WorkerAddressSansLastOctet == cluster.MasterAddressSansLastOctet && cluster.WorkerLastOctet == 0 {
			cluster.WorkerLastOctet = cluster.MasterLastOctet + cluster.NumMaster
		} else {
			cluster.WorkerLastOctet = cluster.MasterLastOctet
		}
		if cluster.WorkerGateway == 0 {
			cluster.WorkerGateway = cluster.MasterGateway
		}
		if cluster.WorkerDomain == "" {
			cluster.WorkerDomain = cluster.MasterDomain
		}

		if cluster.WorkerBaseVmid < 0 {
			return fmt.Errorf("%s - worker_base_vmid should be a positive integer", cluster.Name)
		}
		if cluster.WorkerLastOctet < 0 || cluster.WorkerLastOctet >= 256 {
			return fmt.Errorf("%s - worker_last_octet should be between 1 and 255", cluster.Name)
		}
		if cluster.WorkerGateway < 0 || cluster.WorkerGateway >= 256 {
			return fmt.Errorf("%s - worker_gateway should be between 1 and 255", cluster.Name)
		}

		// Calculate and check IP ranges for masters
		masterRange, err := calculateIPRange(cluster.MasterAddressSansLastOctet, cluster.MasterLastOctet, cluster.NumMaster)
		if err != nil {
			return fmt.Errorf("error calculating master IP range for cluster %s: %v", cluster.Name, err)
		}

		// Calculate and check IP ranges for workers
		workerRange, err := calculateIPRange(cluster.WorkerAddressSansLastOctet, cluster.WorkerLastOctet, cluster.NumWorker)
		if err != nil {
			return fmt.Errorf("error calculating worker IP range for cluster %s: %v", cluster.Name, err)
		}

		// Check for IP range overlaps within the same cluster
		if overlaps(masterRange, workerRange) {
			return fmt.Errorf("IP range overlap detected within cluster %s between master and worker nodes", cluster.Name)
		}

		// Check for IP range overlaps across different clusters
		for _, existingRange := range ipRanges {
			if overlaps(masterRange, existingRange) || overlaps(workerRange, existingRange) {
				return fmt.Errorf("IP range overlap detected in cluster %s", cluster.Name)
			}
		}

		ipRanges = append(ipRanges, masterRange, workerRange)
	}
	return nil
}
