package types

import "net"

// IpRange rappresenta un intervallo di indirizzi IP.
type IpRange struct {
	Start net.IP
	End   net.IP
}
