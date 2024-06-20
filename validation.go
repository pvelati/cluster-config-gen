package main

import (
	"errors"
	"fmt"
	"net"

	"github.com/pvelati/cluster-config-gen/types"
)

// calculateIpRange calculates the range of IPs given a base IP, starting last octet, and count of IPs.
func calculateIpRange(baseIP string, startOctet, count int) (*types.IpRange, error) {
	if net.ParseIP(fmt.Sprintf("%s.1", baseIP)) == nil {
		return nil, fmt.Errorf("invalid base IP: %s", baseIP)
	}

	start := net.ParseIP(fmt.Sprintf("%s.%d", baseIP, startOctet))
	if start == nil {
		return nil, fmt.Errorf("invalid start IP from base: %s and octet: %d", baseIP, startOctet)
	}

	end := net.ParseIP(fmt.Sprintf("%s.%d", baseIP, startOctet+count-1))
	if end == nil {
		return nil, fmt.Errorf("invalid end IP from base: %s and octet: %d", baseIP, startOctet+count-1)
	}

	return &types.IpRange{Start: start, End: end}, nil
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
func overlaps(r1, r2 *types.IpRange) bool {
	return !(ipLessThan(r1.End, r2.Start) || ipLessThan(r2.End, r1.Start))
}

// validateClusters checks the clusters for any configuration issues.
func validateClusters(clusters []types.Cluster) error {
	allErrors := []error{}
	nameSet := map[string]bool{}
	var ipRanges []*types.IpRange

	for _, cluster := range clusters {

		if cluster.Controlplane.Cluster {
			cluster.Controlplane.Num = 4 //3 masters + 1 vip
		} else {
			cluster.Controlplane.Num = 1
		}

		if cluster.Name == "" {
			allErrors = append(allErrors, fmt.Errorf("%s - name cannot be empty", cluster.Name))
		}
		if nameSet[cluster.Name] {
			allErrors = append(allErrors, fmt.Errorf("duplicate cluster name: %s", cluster.Name))
		}
		nameSet[cluster.Name] = true

		if cluster.Compute.Num < 1 {
			allErrors = append(allErrors, fmt.Errorf("%s - the cluster must have at least 1 worker but found %d", cluster.Name, cluster.Compute.Num))
		}
		if cluster.Controlplane.BaseVmid <= 0 {
			allErrors = append(allErrors, fmt.Errorf("%s - master_base_vmid should be a positive integer but found %d", cluster.Name, cluster.Controlplane.BaseVmid))
		}
		if cluster.Controlplane.LastOctet <= 0 || cluster.Controlplane.LastOctet >= 256 {
			allErrors = append(allErrors, fmt.Errorf("%s - master_last_octet should be between 1 and 255 but found %d", cluster.Name, cluster.Controlplane.LastOctet))
		}
		if cluster.Controlplane.GatewayLastOctet <= 0 || cluster.Controlplane.GatewayLastOctet >= 256 {
			allErrors = append(allErrors, fmt.Errorf("%s - master_gateway_last_octet should be between 1 and 255 but found %d", cluster.Name, cluster.Controlplane.GatewayLastOctet))
		}
		if cluster.Controlplane.Domain == "" {
			allErrors = append(allErrors, fmt.Errorf("%s - master_domain cannot be empty", cluster.Name))
		}
		if cluster.Compute.BaseVmid < 0 {
			allErrors = append(allErrors, fmt.Errorf("%s - worker_base_vmid should be a positive integer but found %d", cluster.Name, cluster.Compute.BaseVmid))
		}
		if cluster.Compute.LastOctet < 0 || cluster.Compute.LastOctet >= 256 {
			allErrors = append(allErrors, fmt.Errorf("%s - worker_last_octet should be between 1 and 255 but found %d", cluster.Name, cluster.Compute.LastOctet))
		}
		if cluster.Compute.GatewayLastOctet < 0 || cluster.Compute.GatewayLastOctet >= 256 {
			allErrors = append(allErrors, fmt.Errorf("%s - worker_gateway_last_octet should be between 1 and 255 but found %d", cluster.Name, cluster.Compute.GatewayLastOctet))
		}

		// Calculate and check IP ranges for masters
		masterRange, err := calculateIpRange(cluster.Controlplane.AddressSansLastOctet, cluster.Controlplane.LastOctet, cluster.Controlplane.Num)
		if err != nil {
			allErrors = append(allErrors, fmt.Errorf("error calculating master IP range for cluster %s: %v", cluster.Name, err))
		}

		// Calculate and check IP ranges for workers
		workerRange, err := calculateIpRange(cluster.Compute.AddressSansLastOctet, cluster.Compute.LastOctet, cluster.Compute.Num)
		if err != nil {
			allErrors = append(allErrors, fmt.Errorf("error calculating worker IP range for cluster %s: %v", cluster.Name, err))
		}

		// Check for IP range overlaps within the same cluster
		if overlaps(masterRange, workerRange) {
			allErrors = append(allErrors, fmt.Errorf("IP range overlap detected within cluster %s between master and worker nodes", cluster.Name))
		}

		// Check for IP range overlaps across different clusters
		for _, existingRange := range ipRanges {
			if overlaps(masterRange, existingRange) || overlaps(workerRange, existingRange) {
				allErrors = append(allErrors, fmt.Errorf("IP range overlap detected in cluster %s", cluster.Name))
			}
		}

		ipRanges = append(ipRanges, masterRange, workerRange)
	}
	if len(allErrors) > 0 {
		return errors.Join(allErrors...)
	}
	return nil
}
