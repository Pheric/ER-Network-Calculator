package ipUtils

import (
	"fmt"
	"regexp"
	"strconv"
)

type Ipv4Addr [5]int

// Parses an IPv4 address out of a string. Must not have any protocol or port.
func ParseIpv4(str string) (addr Ipv4Addr, err error) {
	if str == "" {
		return addr, fmt.Errorf("input string is empty")
	}

	// Regex didn't work when I tried to compress it.. so I guess we get to use the expanded version. Written by hand
	ipv4WithCidrRegex := regexp.MustCompile(`(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)(?:\/([1-3]?\d))?`)
	m := ipv4WithCidrRegex.FindAllStringSubmatch(str, -1)
	if m == nil {
		return addr, fmt.Errorf("malformed address")
	}

	matches := m[0]
	if l := len(matches); l < 5 || l > 6 {
		return addr, fmt.Errorf("invalid address by regex")
	}
	matches = matches[1:]
	if matches[4] == "" {
		matches[4] = "-1"
	}

	for i, octet := range matches {
		parsed, err := strconv.Atoi(octet)
		if err != nil {
			return addr, fmt.Errorf("octet #%d is NaN: %s", i+1, octet)
		}

		addr[i] = parsed
	}

	return addr, nil
}

// Returns the padded binary representation of the IP address
func (ip Ipv4Addr) PrintBinary() (s string) {
	for i, octet := range ip[:4] {
		s += fmt.Sprintf("%08b", octet)
		if i < 3 {
			s += "."
		} else if i == 3 && ip.IsCidrFormatted() {
			s += fmt.Sprintf("/%08b", ip.GetPrefix())
		}
	}

	return s
}

// Returns the ip address
func (ip Ipv4Addr) Print() (s string) {
	s = fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
	if ip.IsCidrFormatted() {
		s += fmt.Sprintf("/%d", ip[4])
	}

	return s
}

func (ip Ipv4Addr) IsCidrFormatted() bool {
	return ip[4] != -1
}

func (ip Ipv4Addr) GetPrefix() int {
	return ip[4]
}

// If there's no prefix and the address is a normal private IP address, this will return its network's class
// Otherwise, it prefers the prefix.
// For example, the address '112.17.100.45/16' is a class B; '240.23.18.1' is a class E
func (ip Ipv4Addr) GetClass() int {
	if ip[0] == 10 || ip[0] == 127 || (ip.IsCidrFormatted() && ip[4] == 8 && !(ip[0] >= 240)) {
		return CLASS_A
	} else if ip[0] == 172 && (ip[1] >= 16 && ip[1] <= 31) || (ip.IsCidrFormatted() && ip[4] == 16) {
		return CLASS_B
	} else if ip[0] == 192 && ip[1] == 168 || (ip.IsCidrFormatted() && ip[4] == 24) {
		return CLASS_C
	} else if ip[0] >= 240 {
		return CLASS_E
	} else {
		return NOT_CLASSFUL
	}
}

func (ip Ipv4Addr) GetType() int {
	if ip[0] == 127 {
		return LOOPBACK
	} else if ip[0] == 169 && ip[1] == 254 {
		return APIPA
	} else if ip[0] >= 224 && ip[0] <= 239 {
		return MULTICAST
	} else if ip[3] == 255 {
		return BROADCAST
	} else if ip.Print() == ip.PrintNetworkAddress() {
		return NETWORK
	} else {
		return UNICAST
	}
}

func (ip Ipv4Addr) GetNetmask() (ret [4]int) {
	if !ip.IsCidrFormatted() {
		return ret
	}

	prefix := ip[4]
	for i := 0; i < 4; i++ {
		j := 255
		for ; j > 0 && prefix > 0; j >>= 1 {
			prefix--
		}
		ret[i] = 255 - j
	}
	return ret
}

func (ip Ipv4Addr) PrintNetmask() (s string) {
	netmask := ip.GetNetmask()
	return fmt.Sprintf("%d.%d.%d.%d", netmask[0], netmask[1], netmask[2], netmask[3])
}

func (ip Ipv4Addr) PrintNetworkAddress() (s string) {
	netmask := ip.GetNetmask()

	for i := 0; i < 4; i++ {
		s += strconv.Itoa(ip[i] & netmask[i])
		if i < 3 {
			s += "."
		}
	}
	s += fmt.Sprintf("/%d", ip[4])

	return s
}

// WIP function that auto-subnets an IPv4 address
// @param `nets`: number of desired subnets. Must be <= 65536
// Returns the resultant list of subnets or an error
// @author Racquel Meyer
func (ip Ipv4Addr) Subnet(nets uint) ([]IpAddr, error) {
	// TODO
	return nil, fmt.Errorf("unimplemented")
}