package service

const (
	defaultLoadBalancerSourceRanges = "0.0.0.0/0"
)

// IsAllowAll checks whether the utilnet.IPNet allows traffic from 0.0.0.0/0
func IsAllowAll(ipRange string) bool {

	if ipRange == "0.0.0.0/0" {
		return true
	}
	return false
}
